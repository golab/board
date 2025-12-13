/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package color_test

import (
	"fmt"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/core/color"
)

func TestColor(t *testing.T) {
	assert.Equal(t, color.Black.String(), "B")
	assert.Equal(t, color.White.String(), "W")
	assert.Equal(t, color.Empty.String(), "+")
}

var oppTests = []struct {
	input  color.Color
	output color.Color
}{
	{color.Empty, color.Empty},
	{color.Black, color.White},
	{color.White, color.Black},
}

func TestOpposite(t *testing.T) {
	for i, tt := range oppTests {
		t.Run(fmt.Sprintf("opp%d", i), func(t *testing.T) {
			opp := tt.input.Opposite()
			assert.Equal(t, opp, tt.output)
		})
	}
}

func TestFill(t *testing.T) {
	var fillTests = []struct {
		input  color.Color
		output color.Color
	}{
		{color.Empty, color.FillDame},
		{color.Black, color.FillBlack},
		{color.FillBlack, color.FillBlack},
		{color.White, color.FillWhite},
		{color.FillWhite, color.FillWhite},
		{color.FillDame, color.FillDame},
	}

	for i, tt := range fillTests {
		t.Run(fmt.Sprintf("fill%d", i), func(t *testing.T) {
			fill := tt.input.Fill()
			assert.Equal(t, fill, tt.output)
		})
	}
}

func TestEqual(t *testing.T) {
	var equalTests = []struct {
		col1   color.Color
		col2   color.Color
		output bool
	}{
		{color.Black, color.Black, true},
		{color.White, color.White, true},
		{color.Empty, color.Empty, true},
		{color.Black, color.FillBlack, true},
		{color.FillBlack, color.Black, true},
		{color.FillWhite, color.White, true},
		{color.White, color.FillWhite, true},
		{color.Black, color.White, false},
	}

	for i, tt := range equalTests {
		t.Run(fmt.Sprintf("fill%d", i), func(t *testing.T) {
			e := tt.col1.Equal(tt.col2)
			assert.Equal(t, e, tt.output)
		})
	}
}
