/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package parser

import (
	"testing"

	"github.com/jarednogo/board/internal/assert"
)

var outputTests = []string{
	"(;GM[1])",
	"(;GM[1];B[aa];W[bb](;B[cc];W[dd])(;B[ee];W[ff]))",
	"(;GM[1];C[some comment])",
	"(;GM[1];C[comment \"with\" quotes])",
	"(;GM[1];C[comment [with\\] brackets])",
}

func TestToSGF(t *testing.T) {
	for _, input := range outputTests {
		t.Run(input, func(t *testing.T) {
			p := New(input)
			root, err := p.Parse()
			assert.NoError(t, err)
			output := root.toSGF(true)
			assert.Equal(t, output, input)
		})
	}
}

func TestPass(t *testing.T) {
	sgf := "(;GM[1];B[aa];W[bb];B[tt];W[ss])"
	p := New(sgf)
	root, err := p.Parse()
	if err != nil {
		t.Error(err)
	}
	output := root.toSGF(true)
	if output != "(;GM[1];B[aa];W[bb];B[];W[ss])" {
		t.Errorf("error in reading [tt] pass")
	}
}
