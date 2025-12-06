/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// the core package provides basic functionality to all the major components of the code
package color

// Color is one of Empty, Black, or White
type Color int

const (
	Empty Color = iota
	Black
	White
	FillBlack
	FillWhite
	FillDame
)

// Opposite: Black -> White, White -> Black, Empty -> Empty
func (c Color) Opposite() Color {
	if c == Black {
		return White
	}
	if c == White {
		return Black
	}
	return Empty
}

func (c Color) Fill() Color {
	switch c {
	case Empty:
		return FillDame
	case Black:
		return FillBlack
	case White:
		return FillWhite
	case FillBlack:
		return FillBlack
	case FillWhite:
		return FillWhite
	case FillDame:
		return FillDame
	}
	return FillDame
}

func (c Color) Equal(d Color) bool {
	if c == d {
		return true
	}
	if c == Black && d == FillBlack {
		return true
	}
	if c == FillBlack && d == Black {
		return true
	}
	if c == White && d == FillWhite {
		return true
	}
	if c == FillWhite && d == White {
		return true
	}
	return false
}

// String is just for debugging purposes
func (c Color) String() string {
	if c == Black {
		return "B"
	}
	if c == White {
		return "W"
	}
	return "+"
}
