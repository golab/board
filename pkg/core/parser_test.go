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

	"github.com/jarednogo/board/pkg/core"
)

var fieldTests = []struct {
	input string
	key   string
	value string
}{
	{"(;GM[1])", "GM", "1"},
	{"(;FF[4])", "FF", "4"},
	{"(;CA[UTF-8])", "CA", "UTF-8"},
	{"(;SZ[19])", "SZ", "19"},
	{"(;PB[a b])", "PB", "a b"},
	{"(;C[[1d\\]Player: \"hello world\"])", "C", "[1d]Player: \"hello world\""},
	{"(;W[aa])", "W", "aa"},
	{"(;B[])", "B", ""},
	{"(;GM [1])", "GM", "1"},
}

func TestParser(t *testing.T) {
	for _, tt := range fieldTests {
		t.Run(tt.input, func(t *testing.T) {
			p := core.NewParser(tt.input)
			root, err := p.Parse()
			if err != nil {
				t.Error(err)
				return
			}
			if val, ok := root.Fields[tt.key]; !ok {
				t.Errorf("key not present: %s", tt.key)
			} else if len(val) != 1 {
				t.Errorf("expected length of multifield to be 1, got: %d", len(val))
			} else if val[0] != tt.value {
				t.Errorf("expected value %s, got: %s", tt.value, val[0])
			}
		})
	}
}

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
			p := core.NewParser(input)
			root, err := p.Parse()
			if err != nil {
				t.Error(err)
				return
			}
			output := root.ToSGF(true)
			if output != input {
				t.Errorf("expected %s, got: %s", input, output)
			}
		})
	}
}

var mergeTests = []struct {
	input []string
	num   int
}{
	{[]string{"(;B[aa])", "(;B[bb])"}, 2},
	{[]string{"(;AB[dd])", "(;PB[B];B[qq])", "(;GM[1](;B[aa])(;B[bb]))"}, 4},
}

func TestMerge(t *testing.T) {
	for i, tt := range mergeTests {
		t.Run(fmt.Sprintf("merge%d", i), func(t *testing.T) {
			merged := core.Merge(tt.input)
			p := core.NewParser(merged)
			root, err := p.Parse()
			if err != nil {
				t.Error(err)
				return
			}
			if len(root.Down) != tt.num {
				t.Errorf("expected %d children, got: %d", tt.num, len(root.Down))
				return
			}
		})
	}
}

func TestPass(t *testing.T) {
	sgf := "(;GM[1];B[aa];W[bb];B[tt];W[ss])"
	p := core.NewParser(sgf)
	root, err := p.Parse()
	if err != nil {
		t.Error(err)
	}
	output := root.ToSGF(true)
	if output != "(;GM[1];B[aa];W[bb];B[];W[ss])" {
		t.Errorf("error in reading [tt] pass")
	}
}

func TestEmpty(t *testing.T) {
	sgf := "()"
	p := core.NewParser(sgf)
	_, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

}

func FuzzParser(f *testing.F) {
	testcases := []string{"(;)", "(;GM[1];B[aa];W[bb];B[];W[ss])", "(;GM[1];C[comment \"with\" quotes])"}
	for _, tc := range testcases {
		// add to seed corpus
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		p := core.NewParser(orig)
		// looking for crashes or panics
		_, _ = p.Parse()

	})
}
