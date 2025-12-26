/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package integration

import (
	"math/rand"

	"github.com/jarednogo/board/pkg/core/color"
	"github.com/jarednogo/board/pkg/core/coord"
	"github.com/jarednogo/board/pkg/state"
)

func GenerateRandomSGF(seed int, num int) string {
	r := rand.New(rand.NewSource(int64(seed)))

	size := 19

	s := state.NewState(size)
	numMoves := 0
	col := color.Black
	for numMoves < num {
		x := r.Intn(size)
		y := r.Intn(size)
		c := coord.NewCoord(x, y)
		if s.Board().Legal(c, col) {
			s.AddNode(c, col)
			col = col.Opposite()
			numMoves++
		}
	}
	return s.ToSGF()
}

func GenerateNRandomSGF(seed, num, minLength, maxLength int) []string {
	r := rand.New(rand.NewSource(int64(seed)))
	sgfs := []string{}
	for i := 0; i < num; i++ {
		randLength := r.Intn(maxLength-minLength) + minLength
		sgf := GenerateRandomSGF(i, randLength)
		sgfs = append(sgfs, sgf)
	}
	return sgfs
}
