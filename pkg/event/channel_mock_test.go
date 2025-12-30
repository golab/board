/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package event_test

import (
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/pkg/event"
)

func TestMockEventChannel(t *testing.T) {
	ec := event.NewMockEventChannel()
	assert.Equal(t, ec.GetRoomID(), "")
	ec.SetRoomID("room1")
	assert.Equal(t, ec.GetRoomID(), "room1")

	e := event.NopEvent()
	assert.Equal(t, len(ec.SavedEvents), 0)
	err := ec.SendEvent(e)
	assert.NoError(t, err)
	assert.Equal(t, len(ec.SavedEvents), 1)

	_, err = ec.ReceiveEvent()
	assert.NotNil(t, err)

	ec.QueuedEvents = append(ec.QueuedEvents, e)
	e, err = ec.ReceiveEvent()
	assert.NoError(t, err)
	assert.Equal(t, e.Type(), "nop")

	_, err = ec.ReceiveEvent()
	assert.NotNil(t, err)

	err = ec.Close()
	assert.NoError(t, err)
	assert.Equal(t, ec.Closed, true)

	assert.Equal(t, len(ec.ID()), 36)
}

func TestBlockingMockEventChannel(t *testing.T) {
	ec := event.NewBlockingMockEventChannel()
	go func() {
		ec.OnConnect()
	}()
	<-ec.Ready()

	var err error

	go func() {
		ec.Disconnect()
	}()
	_, err = ec.ReceiveEvent()
	assert.NotNil(t, err)
}

func TestTwoWayMockEventChannel(t *testing.T) {
	ec := event.NewTwoWayMockEventChannel()
	evt := event.NopEvent()
	err := ec.SendEvent(evt)
	assert.NoError(t, err)
	ec.Flush()

	go func() {
		ec.SimulateEvent(evt)
	}()

	evt, err = ec.ReceiveEvent()
	assert.NoError(t, err)

	assert.Equal(t, evt.Type(), "nop")

	ec.Disconnect()
	_, err = ec.ReceiveEvent()
	assert.NotNil(t, err)
}
