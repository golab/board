/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room

import (
	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/state"
)

// helper functions for ogs
// these were more important when ogs was its own package
// still i'm keeping them around in case i split it out again

func (r *Room) HeadColor() core.Color {
	return r.state.Head.Color
}

func (r *Room) PushHead(x, y int, c core.Color) {
	r.state.PushHead(x, y, c)
}

func (r *Room) GenerateFullFrame(t core.TreeJSONType) *core.Frame {
	return r.state.GenerateFullFrame(t)
}

func (r *Room) AddPatternNodes(movesArr []*core.PatternMove) {
	r.state.AddPatternNodes(movesArr)
}

func (r *Room) ToSGF(init bool) string {
	return r.state.ToSGF(init)
}

func (r *Room) CreateStateJSON() *state.StateJSON {
	return r.state.CreateStateJSON()
}
