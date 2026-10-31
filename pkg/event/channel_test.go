/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package event_test

import (
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/pkg/event"
)

type MockReadWriteCloser struct {
	queuedPackets [][]byte
	index         int
}

func NewMockReadWriteCloser(queued [][]byte) *MockReadWriteCloser {
	return &MockReadWriteCloser{
		queuedPackets: queued,
	}
}

func (c *MockReadWriteCloser) Read(arr []byte) (int, error) {
	if c.index >= len(c.queuedPackets) {
		return 0, io.EOF
	}
	packet := c.queuedPackets[c.index]
	m := min(len(arr), len(packet))
	for i := 0; i < m; i++ {
		arr[i] = packet[i]
	}
	c.index++
	return m, nil
}

func (c *MockReadWriteCloser) Write(arr []byte) (int, error) {
	return 0, nil
}

func (c *MockReadWriteCloser) Close() error {
	return nil
}

func TestEventChannel(t *testing.T) {
	queued := [][]byte{
		{61, 0, 0, 0},
		[]byte(`{"event": "test", "value": "somevalue", "userid": "user_123"}`),
	}
	c := NewMockReadWriteCloser(queued)
	ch := event.NewDefaultEventChannel(c)
	defer ch.Close() //nolint:errcheck
	evt, err := ch.ReceiveEvent()
	assert.NoError(t, err)
	assert.Equal(t, evt.Type(), "test")
	assert.Equal(t, evt.Value().(string), "somevalue")
	assert.Equal(t, evt.User(), "user_123")

	err = ch.SendEvent(evt)
	assert.NoError(t, err)
}

func TestEventChannelBig(t *testing.T) {
	msg := ""
	for i := 0; i < 1024; i++ {
		msg += "A"
	}
	e := fmt.Sprintf(`{"event": "test", "value": "%s", "userid": "user_123"}`, msg)
	queued := [][]byte{}
	queued = append(queued, []byte{52, 4, 0, 0})
	for i := 0; i < int(len(e)/64)+1; i++ {
		start := 64 * i
		end := min(64*(i+1), len(e))
		chunk := e[start:end]
		queued = append(queued, []byte(chunk))
	}

	c := NewMockReadWriteCloser(queued)
	ch := event.NewDefaultEventChannel(c)
	defer ch.Close() //nolint:errcheck
	evt, err := ch.ReceiveEvent()
	assert.NoError(t, err)
	_ = evt
	assert.Equal(t, evt.Type(), "test")
	assert.Equal(t, evt.Value().(string), msg)
	assert.Equal(t, evt.User(), "user_123")

	err = ch.SendEvent(evt)
	assert.NoError(t, err)
}

func TestID(t *testing.T) {
	c := NewMockReadWriteCloser([][]byte{})
	ch := event.NewDefaultEventChannel(c)
	assert.Equal(t, len(ch.ID()), 36)
}

// an oversized length is rejected before anything is read or allocated
func TestEventChannelTooLarge(t *testing.T) {
	queued := [][]byte{
		{0xff, 0xff, 0xff, 0xff},
		[]byte(`{"event": "test"}`),
	}
	c := NewMockReadWriteCloser(queued)
	ch := event.NewDefaultEventChannel(c)
	_, err := ch.ReceiveEvent()
	assert.NotNil(t, err)
	assert.Equal(t, c.index, 1)
}

// recordingWriteCloser records every Write and signals it on a channel, so a
// test can wait for asynchronous delivery.
type recordingWriteCloser struct {
	mu      sync.Mutex
	writes  int
	written chan struct{}
}

func newRecordingWriteCloser() *recordingWriteCloser {
	return &recordingWriteCloser{written: make(chan struct{}, 1024)}
}

func (w *recordingWriteCloser) Read([]byte) (int, error) { return 0, io.EOF }

func (w *recordingWriteCloser) Write(p []byte) (int, error) {
	w.mu.Lock()
	w.writes++
	w.mu.Unlock()
	w.written <- struct{}{}
	return len(p), nil
}

func (w *recordingWriteCloser) Close() error { return nil }

func (w *recordingWriteCloser) count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.writes
}

// sends are delivered on the channel's own goroutine, not the caller's
func TestEventChannelDeliversAsync(t *testing.T) {
	w := newRecordingWriteCloser()
	ch := event.NewDefaultEventChannel(w)
	defer ch.Close() //nolint:errcheck

	const n = 5
	for i := 0; i < n; i++ {
		require.NoError(t, ch.SendEvent(event.NewEvent("test", fmt.Sprintf("%d", i))))
	}

	for i := 0; i < n; i++ {
		select {
		case <-w.written:
		case <-time.After(2 * time.Second):
			t.Fatalf("timed out after %d of %d async writes", i, n)
		}
	}
	assert.Equal(t, w.count(), n)
}

// blockingWriteCloser blocks in Write until it is closed, standing in for a
// client that has stopped reading.
type blockingWriteCloser struct {
	closed chan struct{}
	once   sync.Once
}

func newBlockingWriteCloser() *blockingWriteCloser {
	return &blockingWriteCloser{closed: make(chan struct{})}
}

func (w *blockingWriteCloser) Read([]byte) (int, error) {
	<-w.closed
	return 0, io.EOF
}

func (w *blockingWriteCloser) Write(p []byte) (int, error) {
	<-w.closed
	return 0, io.ErrClosedPipe
}

func (w *blockingWriteCloser) Close() error {
	w.once.Do(func() { close(w.closed) })
	return nil
}

// a client that never drains eventually overflows its buffer and is
// disconnected, instead of blocking the sender forever
func TestEventChannelSlowConsumerDisconnects(t *testing.T) {
	w := newBlockingWriteCloser()
	ch := event.NewDefaultEventChannel(w)
	defer ch.Close() //nolint:errcheck

	// far more than the buffer can hold; the writer goroutine is stuck in
	// Write, so the buffer fills and then overflows
	var err error
	for i := 0; i < 100000; i++ {
		if err = ch.SendEvent(event.EmptyEvent()); err != nil {
			break
		}
	}
	require.ErrorIs(t, err, event.ErrSlowConsumer)

	// once overflowed, the connection is closed and further sends are rejected
	assert.ErrorIs(t, ch.SendEvent(event.EmptyEvent()), event.ErrClosed)
}

// a body split across several short reads is put back together
func TestEventChannelSplitReads(t *testing.T) {
	e := `{"event": "test", "value": "abc", "userid": "u"}`
	queued := [][]byte{{byte(len(e)), 0}, {0, 0}}
	for i := 0; i < len(e); i += 5 {
		queued = append(queued, []byte(e[i:min(i+5, len(e))]))
	}
	ch := event.NewDefaultEventChannel(NewMockReadWriteCloser(queued))
	evt, err := ch.ReceiveEvent()
	assert.NoError(t, err)
	assert.Equal(t, evt.Value().(string), "abc")
}

// lockingWriteCloser mimics x/net/websocket: Write and Close share a write
// lock, so Close can't proceed while a Write is stalled. Write only returns
// once release is closed (standing in for a write deadline firing).
type lockingWriteCloser struct {
	wio     sync.Mutex
	release chan struct{}
	entered chan struct{}
	once    sync.Once
}

func newLockingWriteCloser() *lockingWriteCloser {
	return &lockingWriteCloser{
		release: make(chan struct{}),
		entered: make(chan struct{}, 1),
	}
}

func (w *lockingWriteCloser) Read([]byte) (int, error) { return 0, io.EOF }

func (w *lockingWriteCloser) Write(p []byte) (int, error) {
	w.wio.Lock()
	defer w.wio.Unlock()
	select {
	case w.entered <- struct{}{}:
	default:
	}
	<-w.release
	return 0, io.ErrClosedPipe
}

func (w *lockingWriteCloser) Close() error {
	w.wio.Lock()
	defer w.wio.Unlock()
	return nil
}

func (w *lockingWriteCloser) unblock() { w.once.Do(func() { close(w.release) }) }

// overflowing a stalled connection must not block the sender, even when the
// socket's Close waits on the stalled Write
func TestEventChannelSlowConsumerDoesNotBlockSender(t *testing.T) {
	w := newLockingWriteCloser()
	defer w.unblock()
	ch := event.NewDefaultEventChannel(w)

	// get the writer stuck inside Write, holding the write lock
	require.NoError(t, ch.SendEvent(event.EmptyEvent()))
	select {
	case <-w.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("writer never started")
	}

	done := make(chan error, 1)
	go func() {
		var err error
		for i := 0; i < 100000; i++ {
			if err = ch.SendEvent(event.EmptyEvent()); err != nil {
				break
			}
		}
		done <- err
	}()

	select {
	case err := <-done:
		require.ErrorIs(t, err, event.ErrSlowConsumer)
	case <-time.After(2 * time.Second):
		t.Fatal("SendEvent blocked on a stalled connection")
	}
	assert.ErrorIs(t, ch.SendEvent(event.EmptyEvent()), event.ErrClosed)
}
