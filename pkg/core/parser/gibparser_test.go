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

func TestParseGameLine1(t *testing.T) {
	p := NewGIBParser("test line")
	_, err := p.parseGameLine()
	assert.NotNil(t, err)
}

func TestParseGameLine2(t *testing.T) {
	p := NewGIBParser("test line\n")
	s, err := p.parseGameLine()
	require.NoError(t, err)
	assert.Equal(t, s, "test line")
}

func TestParseGameLine3(t *testing.T) {
	p := NewGIBParser("test line\r\n")
	s, err := p.parseGameLine()
	require.NoError(t, err)
	assert.Equal(t, s, "test line")
}

func TestParseGame1(t *testing.T) {
	p := NewGIBParser("\\GS\nline1\nline2\n\\GE\n")
	lines, err := p.parseGame()
	require.NoError(t, err)
	require.Equal(t, len(lines), 2)
	assert.Equal(t, lines[0], "line1")
	assert.Equal(t, lines[1], "line2")
}

func TestParseGame2(t *testing.T) {
	p := NewGIBParser("line1\nline2\n\\GE\n")
	_, err := p.parseGame()
	assert.NotNil(t, err)
}

func TestParseGame3(t *testing.T) {
	p := NewGIBParser("\\GS\nline1\nline2\n")
	_, err := p.parseGame()
	assert.NotNil(t, err)
}

func TestParseValue1(t *testing.T) {
	p := NewGIBParser("foo\\]")
	s, err := p.parseValue()
	require.NoError(t, err)
	assert.Equal(t, s, "foo")
}

func TestParseValue2(t *testing.T) {
	p := NewGIBParser("foo")
	_, err := p.parseValue()
	assert.NotNil(t, err)
}

func TestParseGIBKey1(t *testing.T) {
	p := NewGIBParser("foo=bar")
	s, err := p.parseKey()
	require.NoError(t, err)
	assert.Equal(t, s, "foo")
}

func TestParseGIBKey2(t *testing.T) {
	p := NewGIBParser("foo:bar")
	_, err := p.parseKey()
	assert.NotNil(t, err)
}

func TestParseGIBProperty1(t *testing.T) {
	p := NewGIBParser("\\[foo=bar\\]")
	kv, err := p.parseProperty()
	require.NoError(t, err)
	require.NotNil(t, kv)
	assert.Equal(t, kv.key, "foo")
	assert.Equal(t, kv.value, "bar")
}

func TestParseGIBProperty2(t *testing.T) {
	p := NewGIBParser("foo=bar")
	_, err := p.parseProperty()
	assert.NotNil(t, err)
}

func TestParseGIBProperty3(t *testing.T) {
	p := NewGIBParser("\\[foo=bar")
	_, err := p.parseProperty()
	assert.NotNil(t, err)
}

func TestParseGIBHeader1(t *testing.T) {
	p := NewGIBParser("\\HS\n\\[foo=bar\\]\n\\[baz=bot\\]\n\\HE")
	kvs, err := p.parseHeader()
	require.NoError(t, err)
	require.Equal(t, len(kvs), 2)
	assert.Equal(t, kvs[0].key, "foo")
	assert.Equal(t, kvs[0].value, "bar")
	assert.Equal(t, kvs[1].key, "baz")
	assert.Equal(t, kvs[1].value, "bot")
}

func TestParseGIBHeader2(t *testing.T) {
	p := NewGIBParser("foo")
	_, err := p.parseHeader()
	assert.NotNil(t, err)
}

func TestParseGIBHeader3(t *testing.T) {
	p := NewGIBParser("\\HS\n\\[foo=bar\\]\n")
	_, err := p.parseHeader()
	assert.NotNil(t, err)
}

func TestParseGIB(t *testing.T) {
	p := NewGIBParser("\\HS\n\\[foo=bar\\]\n\\HE\n\\GS\nline1\nline2\n\\GE")
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.Equal(t, len(gibResult.header), 1)
	assert.Equal(t, gibResult.header[0].key, "foo")
	assert.Equal(t, gibResult.header[0].value, "bar")
	require.Equal(t, len(gibResult.game), 2)
	assert.Equal(t, gibResult.game[0], "line1")
	assert.Equal(t, gibResult.game[1], "line2")
}

func TestGameInfoParseProperty1(t *testing.T) {
	p := NewGameInfoParser("foo:bar")
	kv, err := p.parseProperty()
	require.NoError(t, err)
	assert.Equal(t, kv.key, "foo")
	assert.Equal(t, kv.value, "bar")
}

func TestGameInfoParseProperty2(t *testing.T) {
	p := NewGameInfoParser("foo=bar")
	_, err := p.parseProperty()
	assert.NotNil(t, err)
}

func TestGameInfoParseKey(t *testing.T) {
	p := NewGameInfoParser("foo:bar")
	k, err := p.parseKey()
	require.NoError(t, err)
	assert.Equal(t, k, "foo")
}

func TestGameInfoParseValue1(t *testing.T) {
	p := NewGameInfoParser("bar,")
	v := p.parseValue()
	assert.Equal(t, v, "bar")
}

func TestGameInfoParseValue2(t *testing.T) {
	p := NewGameInfoParser("bar")
	v := p.parseValue()
	assert.Equal(t, v, "bar")
}

func TestGameInfoParseValue3(t *testing.T) {
	p := NewGameInfoParser("bar,baz")
	v := p.parseValue()
	assert.Equal(t, v, "bar")
}

func TestGIBToSGF(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:65\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\HE\n\\GS\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := NewGIBParser(data)
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, gibResult)
	node, err := gibResult.ToSGFNode()
	require.NoError(t, err)
	require.NotNil(t, node)
	assert.Equal(t, len(node.toSGF(true)), 80)
}

func TestGIBToSGFRanks(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:65\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\[GAMEBLACKLEVEL=25\\]\n\\[GAMEWHITELEVEL=17\\]\\HE\n\\GS\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := NewGIBParser(data)
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, gibResult)
	node, err := gibResult.ToSGFNode()
	require.NoError(t, err)
	require.NotNil(t, node)
	br := node.GetField("BR")
	wr := node.GetField("WR")
	require.Equal(t, len(br), 1)
	require.Equal(t, len(wr), 1)
	assert.Equal(t, br[0], "8d")
	assert.Equal(t, wr[0], "1k")
}

func TestGIBToSGFDum(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:0,DUM:10\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\HE\n\\GS\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := NewGIBParser(data)
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, gibResult)
	node, err := gibResult.ToSGFNode()
	require.NoError(t, err)
	require.NotNil(t, node)
	km := node.GetField("KM")
	require.Equal(t, len(km), 1)
	assert.Equal(t, km[0], "-10")
}

func TestGIBToSGFHandicap(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:65\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\HE\n\\GS\nINI 0 1 2 &4\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := NewGIBParser(data)
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, gibResult)
	node, err := gibResult.ToSGFNode()
	require.NoError(t, err)
	require.NotNil(t, node)
	require.Equal(t, len(node.GetField("AB")), 2)
}

func TestGRLT0(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:65,GRLT:0,ZIPSU:90\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\HE\n\\GS\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := NewGIBParser(data)
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, gibResult)
	node, err := gibResult.ToSGFNode()
	require.NoError(t, err)
	require.NotNil(t, node)
	re := node.GetField("RE")
	require.Equal(t, len(re), 1)
	assert.Equal(t, re[0], "B+9")
}

func TestGRLT1(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:65,GRLT:1,ZIPSU:90\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\HE\n\\GS\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := NewGIBParser(data)
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, gibResult)
	node, err := gibResult.ToSGFNode()
	require.NoError(t, err)
	require.NotNil(t, node)
	re := node.GetField("RE")
	require.Equal(t, len(re), 1)
	assert.Equal(t, re[0], "W+9")
}

func TestGRLT3(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:65,GRLT:3,ZIPSU:0\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\HE\n\\GS\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := NewGIBParser(data)
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, gibResult)
	node, err := gibResult.ToSGFNode()
	require.NoError(t, err)
	require.NotNil(t, node)
	re := node.GetField("RE")
	require.Equal(t, len(re), 1)
	assert.Equal(t, re[0], "B+R")
}

func TestGRLT4(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:65,GRLT:4,ZIPSU:0\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\HE\n\\GS\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := NewGIBParser(data)
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, gibResult)
	node, err := gibResult.ToSGFNode()
	require.NoError(t, err)
	require.NotNil(t, node)
	re := node.GetField("RE")
	require.Equal(t, len(re), 1)
	assert.Equal(t, re[0], "W+R")
}

func TestGRLT7(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:65,GRLT:7,ZIPSU:0\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\HE\n\\GS\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := NewGIBParser(data)
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, gibResult)
	node, err := gibResult.ToSGFNode()
	require.NoError(t, err)
	require.NotNil(t, node)
	re := node.GetField("RE")
	require.Equal(t, len(re), 1)
	assert.Equal(t, re[0], "B+T")
}

func TestGRLT8(t *testing.T) {
	data := "\\HS\n\\[GAMEINFOMAIN=LINE:19,GONGJE:65,GRLT:8,ZIPSU:0\\]\n\\[WUSERINFO=WNICK:player_white\\]\n\\[BUSERINFO=BNICK:player_black\\]\n\\HE\n\\GS\nSTO 0 2 1 15 15\nSTO 0 3 2 3 15\nSTO 0 4 1 15 3\nSTO 0 5 2 3 3\nSKI 0 6\nSKI 0 7\n\\GE"
	p := NewGIBParser(data)
	gibResult, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, gibResult)
	node, err := gibResult.ToSGFNode()
	require.NoError(t, err)
	require.NotNil(t, node)
	re := node.GetField("RE")
	require.Equal(t, len(re), 1)
	assert.Equal(t, re[0], "W+T")
}

func TestAddHandicap2(t *testing.T) {
	n := &SGFNode{}
	addHandicap(n, 2)
	assert.True(t, n.HasField("AB"))
	ab := n.GetField("AB")
	require.Equal(t, len(ab), 2)
	assert.Equal(t, ab[0], "pd")
	assert.Equal(t, ab[1], "dp")
}

func TestAddHandicap3(t *testing.T) {
	n := &SGFNode{}
	addHandicap(n, 3)
	assert.True(t, n.HasField("AB"))
	ab := n.GetField("AB")
	require.Equal(t, len(ab), 3)
	assert.Equal(t, ab[0], "pd")
	assert.Equal(t, ab[1], "dp")
	assert.Equal(t, ab[2], "pp")
}

func TestAddHandicap4(t *testing.T) {
	n := &SGFNode{}
	addHandicap(n, 4)
	assert.True(t, n.HasField("AB"))
	ab := n.GetField("AB")
	require.Equal(t, len(ab), 4)
	assert.Equal(t, ab[0], "pd")
	assert.Equal(t, ab[1], "dp")
	assert.Equal(t, ab[2], "pp")
	assert.Equal(t, ab[3], "dd")
}

func TestAddHandicap5(t *testing.T) {
	n := &SGFNode{}
	addHandicap(n, 5)
	assert.True(t, n.HasField("AB"))
	ab := n.GetField("AB")
	require.Equal(t, len(ab), 5)
	assert.Equal(t, ab[0], "pd")
	assert.Equal(t, ab[1], "dp")
	assert.Equal(t, ab[2], "pp")
	assert.Equal(t, ab[3], "jj")
	assert.Equal(t, ab[4], "dd")
}

func TestAddHandicap6(t *testing.T) {
	n := &SGFNode{}
	addHandicap(n, 6)
	assert.True(t, n.HasField("AB"))
	ab := n.GetField("AB")
	require.Equal(t, len(ab), 6)
	assert.Equal(t, ab[0], "pd")
	assert.Equal(t, ab[1], "dp")
	assert.Equal(t, ab[2], "pp")
	assert.Equal(t, ab[3], "dd")
	assert.Equal(t, ab[4], "dj")
	assert.Equal(t, ab[5], "pj")
}

func TestAddHandicap7(t *testing.T) {
	n := &SGFNode{}
	addHandicap(n, 7)
	assert.True(t, n.HasField("AB"))
	ab := n.GetField("AB")
	require.Equal(t, len(ab), 7)
	assert.Equal(t, ab[0], "pd")
	assert.Equal(t, ab[1], "dp")
	assert.Equal(t, ab[2], "pp")
	assert.Equal(t, ab[3], "jj")
	assert.Equal(t, ab[4], "dd")
	assert.Equal(t, ab[5], "dj")
	assert.Equal(t, ab[6], "pj")
}

func TestAddHandicap8(t *testing.T) {
	n := &SGFNode{}
	addHandicap(n, 8)
	assert.True(t, n.HasField("AB"))
	ab := n.GetField("AB")
	require.Equal(t, len(ab), 8)
	assert.Equal(t, ab[0], "pd")
	assert.Equal(t, ab[1], "dp")
	assert.Equal(t, ab[2], "pp")
	assert.Equal(t, ab[3], "dd")
	assert.Equal(t, ab[4], "dj")
	assert.Equal(t, ab[5], "pj")
	assert.Equal(t, ab[6], "jd")
	assert.Equal(t, ab[7], "jp")
}

func TestAddHandicap9(t *testing.T) {
	n := &SGFNode{}
	addHandicap(n, 9)
	assert.True(t, n.HasField("AB"))
	ab := n.GetField("AB")
	require.Equal(t, len(ab), 9)
	assert.Equal(t, ab[0], "pd")
	assert.Equal(t, ab[1], "dp")
	assert.Equal(t, ab[2], "pp")
	assert.Equal(t, ab[3], "jj")
	assert.Equal(t, ab[4], "dd")
	assert.Equal(t, ab[5], "dj")
	assert.Equal(t, ab[6], "pj")
	assert.Equal(t, ab[7], "jd")
	assert.Equal(t, ab[8], "jp")
}

func TestMakeRank(t *testing.T) {
	assert.Equal(t, makeRank(35), "9p")
	assert.Equal(t, makeRank(27), "1p")
	assert.Equal(t, makeRank(26), "9d")
	assert.Equal(t, makeRank(18), "1d")
	assert.Equal(t, makeRank(17), "1k")
	assert.Equal(t, makeRank(1), "17k")
	assert.Equal(t, makeRank(0), "18k")
}
