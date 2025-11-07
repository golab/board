/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// the core package provides basic functionality to all the major components of the code
package event

import (
	"encoding/json"
)

type Event interface {
	Type() string
	Value() any
	User() string
	SetType(string)
	SetValue(any)
	SetUser(string)
}

// Event is the basic struct for sending and receiving messages over
// the websockets
type DefaultEvent struct {
	EventType   string `json:"event"`
	EventValue  any    `json:"value"`
	EventUserID string `json:"userid"`
}

func (e *DefaultEvent) Type() string {
	return e.EventType
}

func (e *DefaultEvent) Value() any {
	return e.EventValue
}

func (e *DefaultEvent) User() string {
	return e.EventUserID
}

func (e *DefaultEvent) SetType(t string) {
	e.EventType = t
}

func (e *DefaultEvent) SetValue(v any) {
	e.EventValue = v
}

func (e *DefaultEvent) SetUser(id string) {
	e.EventUserID = id
}

func NewEvent(t string, value any) Event {
	return &DefaultEvent{
		EventType:  t,
		EventValue: value}
}

func EmptyEvent() Event {
	return &DefaultEvent{}
}

// ErrorEvent is a special case of an Event
func ErrorEvent(msg string) Event {
	return NewEvent("error", msg)
}

func FrameEvent(value any) Event {
	return NewEvent("frame", value)
}

// NopEvent signals to the server to do nothing
// (in particular, don't send to clients)
func NopEvent() Event {
	return NewEvent("nop", nil)
}

func EventFromJSON(data []byte) (Event, error) {
	evt := &DefaultEvent{}
	if err := json.Unmarshal(data, evt); err != nil {
		return nil, err
	}
	return evt, nil
}
