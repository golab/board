/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package parser_test

import (
	"fmt"
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/internal/sgfsamples"
	"github.com/golab/board/pkg/core/parser"
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
			p := parser.New(tt.input)
			root, err := p.Parse()
			require.NoError(t, err)
			if val := root.GetField(tt.key); len(val) == 0 {
				t.Errorf("key not present: %s", tt.key)
			} else if len(val) != 1 {
				t.Errorf("expected length of multifield to be 1, got: %d", len(val))
			} else if val[0] != tt.value {
				t.Errorf("expected value %s, got: %s", tt.value, val[0])
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
			merged := parser.Merge(tt.input)
			p := parser.New(merged)
			root, err := p.Parse()
			if err != nil {
				t.Error(err)
				return
			}
			if root.NumChildren() != tt.num {
				t.Errorf("expected %d children, got: %d", tt.num, root.NumChildren())
				return
			}
		})
	}
}

func TestMerge2(t *testing.T) {
	sgf1 := sgfsamples.SimpleFourMoves
	sgf2 := sgfsamples.SimpleEightMoves
	sgf := parser.Merge([]string{sgf1, sgf2})
	p := parser.New(sgf)
	root, err := p.Parse()
	require.NoError(t, err)
	require.Equal(t, root.NumChildren(), 2)
	child1 := root.GetChild(0)
	child2 := root.GetChild(1)
	// both should have PB, PW, and KM as comments
	require.Equal(t, len(child1.GetField("C")), 3)
	require.Equal(t, len(child2.GetField("C")), 3)
}

func TestMerge3(t *testing.T) {
	sgf := parser.Merge([]string{})
	assert.Equal(t, sgf, "")
}

func TestMerge4(t *testing.T) {
	sgf1 := sgfsamples.SimpleFourMoves
	sgf := parser.Merge([]string{sgf1})
	assert.Equal(t, sgf1, sgf)
}

func TestMerge5(t *testing.T) {
	sgf1 := sgfsamples.SimpleFourMoves
	sgf2 := sgfsamples.SimpleEightMoves
	sgf3 := "foobar"
	sgfMerged1 := parser.Merge([]string{sgf1, sgf2, sgf3})
	sgfMerged2 := parser.Merge([]string{sgf1, sgf2})
	assert.Equal(t, sgfMerged1, sgfMerged2)
}

func TestMerge6(t *testing.T) {
	sgf1 := sgfsamples.SimpleFourMoves
	sgf2 := "(;GM[1]SZ[9];B[cc]W[dd])"
	sgf := parser.Merge([]string{sgf1, sgf2})
	assert.Equal(t, sgf, sgf1)
}

func TestEmpty(t *testing.T) {
	sgf := "()"
	p := parser.New(sgf)
	_, err := p.Parse()
	if err != nil {
		t.Error(err)
	}
}

var oddTests = []struct {
	input string
	err   bool
}{
	{"()", false},
	{"(;)", false},
	{"(;;;)", false},
	{"garbage(;GM[1])", false},
	{"(;GM[1])garbageafter", false},
	{"totalgarbage", true},
	{"( ; GM [1] )", false},
	{"garbage (abc, def) stuff", true},
}

func TestOdd(t *testing.T) {
	for _, tt := range oddTests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.New(tt.input)
			_, err := p.Parse()
			assert.Equal(t, err != nil, tt.err)
		})
	}
}

func TestChineseNames(t *testing.T) {
	p := parser.New(sgfsamples.ChineseNames)
	root, err := p.Parse()
	require.NoError(t, err)
	pb := root.GetField("PB")
	pw := root.GetField("PW")

	require.Equal(t, len(pb), 1)
	require.Equal(t, len(pw), 1)

	assert.Equal(t, pw[0], "王思雅")
	assert.Equal(t, pb[0], "李晨宇")
}

func TestMixedCaseField(t *testing.T) {
	p := parser.New(sgfsamples.MixedCaseField)
	root, err := p.Parse()
	require.NoError(t, err)
	c := root.GetField("COPYRIGHT")
	require.Equal(t, len(c), 1)
	assert.Equal(t, c[0], "SomeCopyright")
}

func TestSGFNodeAddField(t *testing.T) {
	n := &parser.SGFNode{}
	n.AddField("foo", "bar")
	n.AddField("baz", "bot")
	assert.Equal(t, len(n.AllFields()), 2)
}

func TestMultifield(t *testing.T) {
	p := parser.New("(;GM[1]ZZ[foo][bar][baz])")
	root, err := p.Parse()
	require.NoError(t, err)
	zz := root.GetField("ZZ")
	require.Equal(t, len(zz), 3)
	assert.Equal(t, zz[0], "foo")
	assert.Equal(t, zz[1], "bar")
	assert.Equal(t, zz[2], "baz")
}

func TestGetChild(t *testing.T) {
	p := parser.New("(;A[b](;C[d])(;E[f]))")
	root, err := p.Parse()
	require.NoError(t, err)
	assert.NotNil(t, root.GetChild(0))
	assert.NotNil(t, root.GetChild(1))
	assert.Zero(t, root.GetChild(2))
	assert.Zero(t, root.GetChild(-1))
}

func TestAnyParse(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:65\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\HE\n\\GS\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := parser.New(data)
	root, err := p.Parse()
	require.NoError(t, err)
	assert.NotNil(t, root)
}

func FuzzParser(f *testing.F) {
	testcases := []string{"(;)", "(;GM[1];B[aa];W[bb];B[];W[ss])", "(;GM[1];C[comment \"with\" quotes])", sgfsamples.Empty, sgfsamples.SimpleTwoBranches, sgfsamples.SimpleWithComment, sgfsamples.SimpleFourMoves, sgfsamples.SimpleEightMoves, sgfsamples.Scoring1, sgfsamples.PassWithTT, sgfsamples.ChineseNames}
	for _, tc := range testcases {
		// add to seed corpus
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		p := parser.New(orig)
		// looking for crashes or panics
		_, _ = p.Parse()
	})
}
