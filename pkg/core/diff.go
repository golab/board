/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// the core package provides basic functionality to all the major components of the code
package core

import (
	"github.com/jarednogo/board/pkg/core/coord"
)

// Diff contains two StoneSets (Add and Remove) and is a key component of a Frame
type Diff struct {
	Add    []*coord.StoneSet `json:"add"`
	Remove []*coord.StoneSet `json:"remove"`
}

// NewDiff makes a Diff based on two StoneSets
func NewDiff(add, remove []*coord.StoneSet) *Diff {
	return &Diff{
		Add:    add,
		Remove: remove,
	}
}

// Copy makes a copy of the Diff
func (d *Diff) Copy() *Diff {
	if d == nil {
		return nil
	}
	add := []*coord.StoneSet{}
	remove := []*coord.StoneSet{}
	for _, a := range d.Add {
		add = append(add, a.Copy())
	}
	for _, r := range d.Remove {
		remove = append(remove, r.Copy())
	}
	return NewDiff(add, remove)
}

// Invert simply exchanges Add and Remove
func (d *Diff) Invert() *Diff {
	if d == nil {
		return nil
	}
	return NewDiff(d.Remove, d.Add)
}
