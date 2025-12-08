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
	"github.com/jarednogo/board/internal/require"
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
			require.NoError(t, err)
			output := root.toSGF(true)
			assert.Equal(t, output, input)
		})
	}
}

func TestPass(t *testing.T) {
	sgf := "(;GM[1];B[aa];W[bb];B[tt];W[ss])"
	p := New(sgf)
	root, err := p.Parse()
	require.NoError(t, err)
	output := root.toSGF(true)
	if output != "(;GM[1];B[aa];W[bb];B[];W[ss])" {
		t.Errorf("error in reading [tt] pass")
	}
}

func TestParseKey1(t *testing.T) {
	text := ""
	p := New(text)
	_, err := p.parseKey()
	assert.NotNil(t, err)
}

func TestParseKey2(t *testing.T) {
	text := "ABC[123]"
	p := New(text)
	key, err := p.parseKey()
	require.NoError(t, err)
	assert.Equal(t, key, "ABC")
}

func TestParseKey3(t *testing.T) {
	text := "Abc[123]"
	p := New(text)
	key, err := p.parseKey()
	require.NoError(t, err)
	assert.Equal(t, key, "ABC")
}

func TestParseKey4(t *testing.T) {
	text := "FOO"
	p := New(text)
	_, err := p.parseKey()
	assert.NotNil(t, err)
}

func TestParseField1(t *testing.T) {
	text := "[123]"
	p := New(text)
	field, err := p.parseField()
	require.NoError(t, err)
	assert.Equal(t, field, "123")
}

func TestParseField2(t *testing.T) {
	text := "123"
	p := New(text)
	_, err := p.parseField()
	assert.NotNil(t, err)
}

func TestParseField3(t *testing.T) {
	text := "[123\\]abc]"
	p := New(text)
	field, err := p.parseField()
	require.NoError(t, err)
	assert.Equal(t, field, "123]abc")
}

func TestParseField4(t *testing.T) {
	text := "[abc"
	p := New(text)
	_, err := p.parseField()
	assert.NotNil(t, err)
}

func TestSkipUntil(t *testing.T) {
	var skipTests = []struct {
		input    string
		output   string
		hasError bool
	}{
		{"garbagedata()", "garbagedata", false},
		{"garbage", "", true},
	}

	for _, tt := range skipTests {
		t.Run(tt.input, func(t *testing.T) {
			p := New(tt.input)
			s, err := p.skipUntil('(')
			if tt.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, s, tt.output)
			}
		})
	}
}

func TestParseNode1(t *testing.T) {
	text := ";FOO[bar][baz][bot];"
	p := New(text)
	node, err := p.parseNode()
	require.NoError(t, err)
	require.Equal(t, len(node.AllFields()), 1)
	g := node.GetField("FOO")
	require.Equal(t, len(g), 3)
	assert.Equal(t, g[0], "bar")
	assert.Equal(t, g[1], "baz")
	assert.Equal(t, g[2], "bot")
}

func TestParseNode2(t *testing.T) {
	text := ";FOO[bar][baz][bot]QUX[quz];"
	p := New(text)
	node, err := p.parseNode()
	require.NoError(t, err)
	require.Equal(t, len(node.AllFields()), 2)
	g := node.GetField("FOO")
	require.Equal(t, len(g), 3)
	assert.Equal(t, g[0], "bar")
	assert.Equal(t, g[1], "baz")
	assert.Equal(t, g[2], "bot")
	h := node.GetField("QUX")
	require.Equal(t, len(h), 1)
	assert.Equal(t, h[0], "quz")
}

func TestParseNode3(t *testing.T) {
	text := "FOO[bb];"
	p := New(text)
	_, err := p.parseNode()
	assert.NotNil(t, err)
}

func TestParseNode4(t *testing.T) {
	text := ";FOO[bb]BAR;"
	p := New(text)
	_, err := p.parseNode()
	assert.NotNil(t, err)
}

func TestParseNodes1(t *testing.T) {
	text := ";FOO[bar"
	p := New(text)
	_, err := p.parseOneOrMoreNodes()
	assert.NotNil(t, err)
}

func TestParseNodes2(t *testing.T) {
	text := ";FOO[bar];BAZ[bot"
	p := New(text)
	_, err := p.parseOneOrMoreNodes()
	assert.NotNil(t, err)
}

func TestParseProperty1(t *testing.T) {
	text := "789[abc]"
	p := New(text)
	_, err := p.parseProperty()
	assert.NotNil(t, err)
}

func TestParseProperty2(t *testing.T) {
	// the SGF spec says keys should be all uppercase
	// so our parser is a bit permissive
	text := "foo[abc]"
	p := New(text)
	prop, err := p.parseProperty()
	require.NoError(t, err)
	assert.Equal(t, prop.key, "FOO")
	require.Equal(t, len(prop.values), 1)
	assert.Equal(t, prop.values[0], "abc")
}

func TestParseProperty3(t *testing.T) {
	// perhaps the parser is TOO permissive
	text := "F O O  [abc]"
	p := New(text)
	prop, err := p.parseProperty()
	require.NoError(t, err)
	assert.Equal(t, prop.key, "FOO")
	require.Equal(t, len(prop.values), 1)
	assert.Equal(t, prop.values[0], "abc")
}

func TestParseProperty4(t *testing.T) {
	text := "FOO"
	p := New(text)
	_, err := p.parseProperty()
	assert.NotNil(t, err)
}

func TestParseProperty5(t *testing.T) {
	text := "FOO[bar][baz"
	p := New(text)
	_, err := p.parseProperty()
	assert.NotNil(t, err)
}

func TestOneOrMoreFields1(t *testing.T) {
	text := "abc"
	p := New(text)
	_, err := p.parseOneOrMoreFields("FOO")
	assert.NotNil(t, err)
}

func TestOneOrMoreFields2(t *testing.T) {
	text := "[abc][def"
	p := New(text)
	_, err := p.parseOneOrMoreFields("FOO")
	assert.NotNil(t, err)
}

func TestParseOneField1(t *testing.T) {
	text := "[abc"
	p := New(text)
	_, err := p.parseOneField("FOO")
	assert.NotNil(t, err)
}

func TestParseOneField2(t *testing.T) {
	text := "[tt]"
	p := New(text)
	s, err := p.parseOneField("B")
	require.NoError(t, err)
	assert.Equal(t, s, "")
}

func TestParseOneField3(t *testing.T) {
	text := "[bb]"
	p := New(text)
	s, err := p.parseOneField("B")
	require.NoError(t, err)
	assert.Equal(t, s, "bb")
}

func TestParseBranch1(t *testing.T) {
	text := "abc"
	p := New(text)
	_, err := p.parseBranch()
	assert.NotNil(t, err)
}

func TestParseBranch2(t *testing.T) {
	text := "("
	p := New(text)
	_, err := p.parseBranch()
	assert.NotNil(t, err)
}

func TestParseBranch3(t *testing.T) {
	text := "(;GM[1](B[aa])(;B[bb]))"
	p := New(text)
	_, err := p.parseBranch()
	assert.NotNil(t, err)
}

func TestParseBranch4(t *testing.T) {
	text := "((;GM[1](;B[aa])(;B[bb])))"
	p := New(text)
	_, err := p.parseBranch()
	require.NoError(t, err)
}

func TestParseBranch5(t *testing.T) {
	text := "(;[1])"
	p := New(text)
	_, err := p.parseBranch()
	assert.NotNil(t, err)
}
