/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package plugin

import (
	"github.com/jarednogo/board/pkg/core/color"
	"github.com/jarednogo/board/pkg/core/coord"
	"github.com/jarednogo/board/pkg/event"
)

type MockPlugin struct {
	IsStarted bool
}

func NewMockPlugin() *MockPlugin {
	return &MockPlugin{}
}

func (mp *MockPlugin) Start(map[string]any) {
	mp.IsStarted = true
}

func (mp *MockPlugin) End() {
	mp.IsStarted = false
}

type call struct {
	name string
	args []any
}

type MockRoom struct {
	calls []*call
}

func (r *MockRoom) HeadColor() color.Color {
	r.calls = append(r.calls, &call{name: "HeadColor"})
	return color.Black
}

func (r *MockRoom) PushHead(x, y int, c color.Color) bool {
	r.calls = append(r.calls, &call{name: "PushHead", args: []any{x, y, c}})
	return true
}

func (r *MockRoom) BroadcastFullFrame() {
	r.calls = append(r.calls, &call{name: "BroadcastFullFrame"})
}

func (r *MockRoom) BroadcastTreeOnly() {
	r.calls = append(r.calls, &call{name: "BroadcastTreeOnly"})
}

func (r *MockRoom) AddStonesToTrunk(t int, s []*coord.Stone) {
	r.calls = append(r.calls, &call{name: "AddStonesToTrunk", args: []any{t, s}})
}

func (r *MockRoom) GetColorAt(t int) color.Color {
	r.calls = append(r.calls, &call{name: "GetColorAt", args: []any{t}})
	if t%2 == 0 {
		return color.White
	}
	return color.Black
}

func (r *MockRoom) Broadcast(e event.Event) {
	r.calls = append(r.calls, &call{name: "Broadcast", args: []any{e}})
}

func (r *MockRoom) UploadSGF(s string) event.Event {
	r.calls = append(r.calls, &call{name: "UploadSGF", args: []any{s}})
	return event.NewEvent("testevent", nil)
}
