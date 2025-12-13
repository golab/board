/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"github.com/jarednogo/board/pkg/core/color"
	"github.com/jarednogo/board/pkg/core/coord"
)

type FieldProvider interface {
	GetField(string) []string
}

func IsMove(f FieldProvider) bool {
	bvalues := f.GetField("B")
	wvalues := f.GetField("W")
	return len(bvalues) > 0 || len(wvalues) > 0
}

func IsPass(f FieldProvider) bool {
	bvalues := f.GetField("B")
	wvalues := f.GetField("W")
	return (len(bvalues) == 1 && bvalues[0] == "") ||
		(len(wvalues) == 1 && wvalues[0] == "")
}

func Color(f FieldProvider) color.Color {
	bvalues := f.GetField("B")
	wvalues := f.GetField("W")
	if len(bvalues) > 0 {
		return color.Black
	}
	if len(wvalues) > 0 {
		return color.White
	}
	return color.Empty
}

func Coord(f FieldProvider) *coord.Coord {
	bvalues := f.GetField("B")
	wvalues := f.GetField("W")
	if len(bvalues) == 1 {
		return coord.FromLetters(bvalues[0])
	}
	if len(wvalues) == 1 {
		return coord.FromLetters(wvalues[0])
	}
	return nil
}

func HasComment(f FieldProvider) bool {
	comments := f.GetField("C")
	return len(comments) > 0
}
