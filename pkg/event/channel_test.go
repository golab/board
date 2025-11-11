/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package event_test

import (
	"io"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/event"
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
	evt, err := ch.ReceiveEvent()
	assert.NoError(t, err)
	assert.Equal(t, evt.Type(), "test")
	assert.Equal(t, evt.Value().(string), "somevalue")
	assert.Equal(t, evt.User(), "user_123")

	err = ch.SendEvent(evt)
	assert.NoError(t, err)
}
