/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package event

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// MaxMessageBytes caps one incoming message; the largest legitimate
	// ones are uploaded files (base64 encoded)
	MaxMessageBytes = 16 << 20

	// messageTimeout is how long the body of a message can take to arrive
	messageTimeout = 30 * time.Second

	// sendQueueSize is the per-connection outbound buffer. Each connection
	// drains its own queue on a dedicated goroutine, so a slow client can
	// fall this far behind before it is disconnected without blocking any
	// other client or the room lock.
	sendQueueSize = 256
)

var (
	// ErrSlowConsumer is returned by SendEvent when the connection's outbound
	// buffer is full. The connection is closed so a stalled client can't hold
	// up broadcasts to everyone else.
	ErrSlowConsumer = errors.New("slow consumer: send buffer full")

	// ErrClosed is returned by SendEvent after the connection has been closed.
	ErrClosed = errors.New("event channel closed")
)

type EventChannel interface {
	SendEvent(evt Event) error
	ReceiveEvent() (Event, error)
	OnConnect()
	Close() error
	ID() string
}

// DefaultEventChannel wraps a ReadWriteCloser. Sends are buffered and written
// on a dedicated per-connection goroutine so that a single slow client can
// never stall a broadcast to everyone else (or the room lock the broadcast is
// holding). Reads stay synchronous on the caller's goroutine.
type DefaultEventChannel struct {
	ws       io.ReadWriteCloser
	id       string
	queue    chan Event
	done     chan struct{}
	start    sync.Once
	closeOne sync.Once
	closeErr error
}

func NewDefaultEventChannel(ws io.ReadWriteCloser) EventChannel {
	// assign the new connection a new id
	id := uuid.New().String()

	return &DefaultEventChannel{
		ws:    ws,
		id:    id,
		queue: make(chan Event, sendQueueSize),
		done:  make(chan struct{}),
	}
}

func (ec *DefaultEventChannel) ID() string {
	return ec.id
}

// SendEvent enqueues evt for delivery and returns without touching the socket.
// It never blocks: if the connection has been closed it returns ErrClosed, and
// if the outbound buffer is full (a stalled client) it closes the connection
// and returns ErrSlowConsumer rather than blocking the caller.
func (ec *DefaultEventChannel) SendEvent(evt Event) error {
	select {
	case <-ec.done:
		return ErrClosed
	default:
	}

	// start the writer lazily, so a channel that never sends never spawns a
	// goroutine
	ec.start.Do(func() { go ec.sendLoop() })

	select {
	case ec.queue <- evt:
		return nil
	default:
		_ = ec.Close()
		return ErrSlowConsumer
	}
}

// sendLoop is the sole writer to the socket: it drains the outbound queue in
// order until the connection is closed or a write fails.
func (ec *DefaultEventChannel) sendLoop() {
	for {
		select {
		case <-ec.done:
			return
		case evt := <-ec.queue:
			data, err := json.Marshal(evt)
			if err != nil {
				// a single unmarshalable event shouldn't kill the connection
				continue
			}
			if _, err := ec.ws.Write(data); err != nil {
				_ = ec.Close()
				return
			}
		}
	}
}

func (ec *DefaultEventChannel) OnConnect() {
}

func (ec *DefaultEventChannel) ReceiveEvent() (Event, error) {
	data, err := ec.readPacket()
	if err != nil {
		return nil, err
	}

	// turn data into json
	var evt Event
	if evt, err = EventFromJSON(data); err != nil {
		return nil, fmt.Errorf("error: %w, string: \"%s\"", err, string(data))
	}
	return evt, nil
}

// Close stops the writer goroutine and closes the underlying socket. It is
// idempotent and safe to call from the sender goroutine, a broadcaster, or the
// connection handler.
func (ec *DefaultEventChannel) Close() error {
	ec.closeOne.Do(func() {
		close(ec.done)
		ec.closeErr = ec.ws.Close()
	})
	return ec.closeErr
}

func (ec *DefaultEventChannel) readPacket() ([]byte, error) {
	// read in 4 bytes (length of rest of message)
	lengthArray := make([]byte, 4)
	_, err := io.ReadFull(ec.ws, lengthArray)
	if err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lengthArray)

	// the length comes from the client, so check it before trusting it
	if length > MaxMessageBytes {
		return nil, fmt.Errorf("message too large: %d bytes (max %d)", length, MaxMessageBytes)
	}

	// once a message has started, the rest of it has to arrive promptly
	// (waiting between messages is fine, so this is only set here)
	if d, ok := ec.ws.(interface{ SetReadDeadline(time.Time) error }); ok {
		_ = d.SetReadDeadline(time.Now().Add(messageTimeout))
		defer func() { _ = d.SetReadDeadline(time.Time{}) }()
	}

	// read in the rest of the data
	data := make([]byte, length)
	_, err = io.ReadFull(ec.ws, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
