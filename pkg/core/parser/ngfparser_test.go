/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package parser

import (
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
)

func TestNGFParseLine(t *testing.T) {
	s := "some title\n"
	p := NewNGFParser(s)
	parsed, err := p.parseLine()
	require.NoError(t, err)
	assert.Equal(t, parsed, "some title")
}

func TestNGFParseInt(t *testing.T) {
	s := "19\n"
	p := NewNGFParser(s)
	n, err := p.parseInt()
	require.NoError(t, err)
	assert.Equal(t, n, 19)
}

func TestNGFParseNickRank1(t *testing.T) {
	s := "user1 10k\n"
	p := NewNGFParser(s)
	data, err := p.parseNickRank()
	require.NoError(t, err)
	assert.Equal(t, data.nick, "user1")
	assert.Equal(t, data.rank, "10k")
}

func TestNGFParseNickRank2(t *testing.T) {
	s := "user1 10k*\n"
	p := NewNGFParser(s)
	data, err := p.parseNickRank()
	require.NoError(t, err)
	assert.Equal(t, data.nick, "user1")
	assert.Equal(t, data.rank, "10k*")
}

func TestNGFParseNickRank3(t *testing.T) {
	s := "user 1 10k\n"
	p := NewNGFParser(s)
	data, err := p.parseNickRank()
	require.NoError(t, err)
	assert.Equal(t, data.nick, "user 1")
	assert.Equal(t, data.rank, "10k")
}

func TestNGFParseNickRank4(t *testing.T) {
	s := "user 1 1 10k\n"
	p := NewNGFParser(s)
	data, err := p.parseNickRank()
	require.NoError(t, err)
	assert.Equal(t, data.nick, "user 1 1")
	assert.Equal(t, data.rank, "10k")
}

func TestNGFParseResult1(t *testing.T) {
	s := "White wins by resignation"
	p := NewNGFParser(s)
	result, err := p.parseResult()
	require.NoError(t, err)
	assert.Equal(t, result, "W+R")
}

func TestNGFParseResult2(t *testing.T) {
	s := "Black wins by resignation"
	p := NewNGFParser(s)
	result, err := p.parseResult()
	require.NoError(t, err)
	assert.Equal(t, result, "B+R")
}

func TestNGFParseResult3(t *testing.T) {
	s := "White wins by time"
	p := NewNGFParser(s)
	result, err := p.parseResult()
	require.NoError(t, err)
	assert.Equal(t, result, "W+T")
}

func TestNGFParseResult4(t *testing.T) {
	s := "Black wins by time"
	p := NewNGFParser(s)
	result, err := p.parseResult()
	require.NoError(t, err)
	assert.Equal(t, result, "B+T")
}

func TestNGFParseResult5(t *testing.T) {
	s := "Black wins by 16.5 points"
	p := NewNGFParser(s)
	result, err := p.parseResult()
	require.NoError(t, err)
	assert.Equal(t, result, "B+16.5")
}

func TestNGFParseResult6(t *testing.T) {
	s := "White wins by 0.5 points"
	p := NewNGFParser(s)
	result, err := p.parseResult()
	require.NoError(t, err)
	assert.Equal(t, result, "W+0.5")
}

func TestNGFParseResult7(t *testing.T) {
	s := "Progressing"
	p := NewNGFParser(s)
	result, err := p.parseResult()
	require.NoError(t, err)
	assert.Equal(t, result, "")
}

func TestNGFParseResult8(t *testing.T) {
	s := "Black wins points"
	p := NewNGFParser(s)
	result, err := p.parseResult()
	require.NoError(t, err)
	assert.Equal(t, result, "B+")
}

func TestNGFMove(t *testing.T) {
	s := "PMABWQRRQ\n"
	p := NewNGFParser(s)
	mv, err := p.parseMove()
	require.NoError(t, err)
	assert.Equal(t, mv.key, "W")
	assert.Equal(t, mv.coord, "pq")
}

const testNGFString = `
SomeApp Friendly Match
19
w1234      5K*
b5678      5K
someapp.com
0
0
6
20220222 [02:22]
0
White loses on time
7
PMABBREER
PMACWEQQE
PMADBQRRQ
PMAEWEEEE
PMAFBPDDP
PMJMWDGGD
PMJNBEKKE`

func TestNGF1(t *testing.T) {
	p := NewNGFParser(testNGFString)
	r, err := p.Parse()
	require.NoError(t, err)
	assert.Equal(t, r.title, "SomeApp Friendly Match")
	assert.Equal(t, r.size, 19)

	require.NotNil(t, r.userWhite)
	assert.Equal(t, r.userWhite.nick, "w1234")
	assert.Equal(t, r.userWhite.rank, "5K*")

	require.NotNil(t, r.userBlack)
	assert.Equal(t, r.userBlack.nick, "b5678")
	assert.Equal(t, r.userBlack.rank, "5K")

	assert.Equal(t, r.handicap, 0)

	assert.Equal(t, r.komi, 6.5)

	assert.Equal(t, r.date, "20220222 [02:22]")

	assert.Equal(t, r.result, "W+T")

	assert.Equal(t, len(r.moves), 7)
}

func TestNGF2(t *testing.T) {
	p := NewNGFParser(testNGFString)
	r, err := p.Parse()
	require.NoError(t, err)
	node, err := r.ToSGFNode()
	require.NoError(t, err)
	t.Logf("%v", node.toSGF(true))
}
