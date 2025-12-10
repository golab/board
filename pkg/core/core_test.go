/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core_test

import (
	"fmt"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/internal/sgfsamples"
	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/core/color"
	"github.com/jarednogo/board/pkg/core/coord"
	"github.com/jarednogo/board/pkg/core/fields"
	"github.com/jarednogo/board/pkg/state"
)

var sanitizeTests = []struct {
	input  string
	output string
}{
	{"  AAA  ", "aaa"},
	{"http://www.example.com?foo=Bar&baz=BOT", "httpwwwexamplecomfoobarbazbot"},
	{"AbC-123", "abc-123"},
	{"MyCoolBoard", "mycoolboard"},
}

func TestSanitize(t *testing.T) {
	for i, tt := range sanitizeTests {
		t.Run(fmt.Sprintf("sanitize%d", i), func(t *testing.T) {
			output := core.Sanitize(tt.input)
			assert.Equal(t, tt.output, output)
		})
	}
}

func TestUUID(t *testing.T) {
	s := core.UUID4()
	assert.Equal(t, len(s), 36)
}

// be a bit wary. this is a non-deterministic test
func TestRandomBoardName(t *testing.T) {
	count := make(map[string]int)
	for i := 0; i < 1000; i++ {
		name := core.RandomBoardName()
		if _, ok := count[name]; !ok {
			count[name] = 0
		}
		count[name] += 1
	}
	for _, v := range count {
		assert.Equal(t, v, 1)
	}
}

func TestIsMove1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("B", "aa")
	assert.True(t, core.IsMove(f))
}

func TestIsMove2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("W", "aa")
	assert.True(t, core.IsMove(f))
}

func TestIsMove3(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("b", "aa")
	assert.False(t, core.IsMove(f))
}

func TestIsMove4(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("w", "aa")
	assert.False(t, core.IsMove(f))
}

func TestIsPass1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("B", "")
	assert.True(t, core.IsPass(f))
}

func TestIsPass2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("W", "")
	assert.True(t, core.IsPass(f))
}

func TestColor1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("B", "aa")
	assert.Equal(t, core.Color(f), color.Black)
}

func TestColor2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("W", "aa")
	assert.Equal(t, core.Color(f), color.White)
}

func TestColor3(t *testing.T) {
	f := &fields.Fields{}
	assert.Equal(t, core.Color(f), color.Empty)
}

func TestCoord1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("B", "aa")
	crd := core.Coord(f)
	assert.True(t, crd.Equal(coord.NewCoord(0, 0)))
}

func TestCoord2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("W", "bb")
	crd := core.Coord(f)
	assert.True(t, crd.Equal(coord.NewCoord(1, 1)))
}

func TestCoordNil(t *testing.T) {
	f := &fields.Fields{}
	crd := core.Coord(f)
	assert.Zero(t, crd)
}

func TestTreeNodeIsMove(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.Resignation1)
	assert.NoError(t, err)

	root := s.Root()
	assert.False(t, core.IsMove(root))

	assert.True(t, len(root.Down) > 0)
	current := root.Down[0]
	assert.True(t, core.IsMove(current))

	assert.True(t, len(current.Down) > 0)
	current = current.Down[0]
	assert.True(t, core.IsMove(current))
}
