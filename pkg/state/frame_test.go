/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
)

// regression
func TestLabel(t *testing.T) {
	s := NewState(19)
	s.current.AddField("LB", "dd:>:D")

	marks := s.generateMarks()
	require.Equal(t, len(marks.Labels), 1)
	assert.Equal(t, marks.Labels[0].Text, ">:D")
}

func TestGenerateTreeOnly(t *testing.T) {
	s := NewState(19)
	f := s.GenerateTreeOnly(Full)
	assert.Equal(t, f.Type, DiffFrame)
}

// frames are serialized later on per-connection goroutines, outside the room
// lock, so they must not share memory with the live state
func TestFullFrameDoesNotAliasState(t *testing.T) {
	s := NewState(19)
	s.root.OverwriteField("ZZ", "z")
	s.root.AppendField("AA", []string{"a", "b"})
	s.current.AppendField("C", []string{"one", "two"})

	f := s.GenerateFullFrame(Full)
	wantKeys := []string{}
	for _, fld := range s.root.AllFields() {
		wantKeys = append(wantKeys, fld.Key)
	}

	// in-place mutations: RemoveField shifts the backing arrays, SortFields
	// reorders the field slice
	s.root.RemoveField("AA", "a")
	s.root.SortFields()
	s.current.RemoveField("C", "one")

	require.Equal(t, len(f.Comments), 2)
	assert.Equal(t, f.Comments[0], "one")
	assert.Equal(t, f.Comments[1], "two")

	var aa []string
	keys := []string{}
	for _, fld := range f.Metadata.Fields {
		keys = append(keys, fld.Key)
		if fld.Key == "AA" {
			aa = fld.Values
		}
	}
	require.Equal(t, len(aa), 2)
	assert.Equal(t, aa[0], "a")
	assert.Equal(t, aa[1], "b")
	// original (unsorted) order is preserved in the frame
	require.Equal(t, len(keys), len(wantKeys))
	for i := range keys {
		assert.Equal(t, keys[i], wantKeys[i])
	}
}
