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

func TestEvent(t *testing.T) {
	e := event.NewEvent("test", nil)

	e.SetType("test123")
	assert.Equal(t, e.Type(), "test123")

	e.SetValue(42)
	assert.Equal(t, e.Value().(int), 42)

	e.SetUser("user123")
	assert.Equal(t, e.User(), "user123")
}

func TestEventCreators(t *testing.T) {
	e := event.EmptyEvent()
	assert.Equal(t, e.Type(), "")

	e = event.ErrorEvent("test msg")
	assert.Equal(t, e.Type(), "error")

	e = event.FrameEvent(nil)
	assert.Equal(t, e.Type(), "frame")

	e = event.NopEvent()
	assert.Equal(t, e.Type(), "nop")
}

func TestDecodeEvent(t *testing.T) {
	data := []byte(`{"event": "test", "value": "somevalue", "userid": "user123"}`)
	evt, err := event.EventFromJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, evt.Type(), "test")
}

func TestDecodeEventFail(t *testing.T) {
	data := []byte("data123")
	_, err := event.EventFromJSON(data)
	assert.NotNil(t, err)
}
