/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package parser

import (
	"fmt"
	"strings"

	"github.com/golab/board/pkg/core/fields"
)

func isWhitespace(c rune) bool {
	return c == '\n' || c == ' ' || c == '\t' || c == '\r'
}

func isUpper(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c rune) bool {
	return c >= 'a' && c <= 'z'
}

type SGFNode struct {
	fields.Fields
	down []*SGFNode
}

// GetChild and NumChildren aren't just for parser_test
// also used in state

func (n *SGFNode) GetChild(i int) *SGFNode {
	if i >= 0 && i < len(n.down) {
		return n.down[i]
	}
	return nil
}

func (n *SGFNode) NumChildren() int {
	return len(n.down)
}

func (n *SGFNode) toSGF(root bool) string {
	sb := strings.Builder{}
	if root {
		sb.WriteByte('(')
	}
	sb.WriteByte(';')

	n.SortFields()

	for _, field := range n.AllFields() {
		sb.WriteString(field.Key)
		for _, value := range field.Values {
			sb.WriteByte('[')
			m := strings.ReplaceAll(value, "]", "\\]")
			sb.WriteString(m)
			sb.WriteByte(']')
		}
	}

	for _, d := range n.down {
		if len(n.down) > 1 {
			sb.WriteByte('(')
			sb.WriteString(d.toSGF(false))
			sb.WriteByte(')')
		} else {
			sb.WriteString(d.toSGF(false))
		}
	}
	if root {
		sb.WriteByte(')')
	}
	return sb.String()
}

type Parser struct {
	*BaseParser
}

func New(text string) *Parser {
	return &Parser{&BaseParser{[]rune(text), 0}}
}

func (p *Parser) isGIB() bool {
	return p.peek(0) == '\\' && p.peek(1) == 'H' && p.peek(2) == 'S'
}

func (p *Parser) isNGF() bool {
	savedIndex := p.index
	defer func() { p.index = savedIndex }()
	p.parseLine()
	_, err := p.parseInt()
	return err == nil
}

func (p *Parser) Parse() (*SGFNode, error) {
	p.skipWhitespace()
	savedIndex := p.index

	// perhaps gib?
	if p.isGIB() {
		gibResult, err := NewGIBParser(string(p.text)).Parse()
		if err == nil {
			return gibResult.ToSGFNode()
		}
	}
	p.index = savedIndex

	// perhaps ngf?
	if p.isNGF() {
		ngfResult, err := NewNGFParser(string(p.text)).Parse()
		if err == nil {
			return ngfResult.ToSGFNode()
		}
	}
	p.index = savedIndex

	// now parse sgf
	return NewSGFParser(string(p.text)).Parse()
}

func Merge(sgfs []string) string {
	if len(sgfs) == 0 {
		return ""
	} else if len(sgfs) == 1 {
		return sgfs[0]
	}

	newRoot := &SGFNode{}

	newRoot.AddField("GM", "1")
	newRoot.AddField("FF", "4")
	newRoot.AddField("CA", "UTF-8")
	newRoot.AddField("PB", "Black")
	newRoot.AddField("PW", "White")
	newRoot.AddField("RU", "Japanese")
	newRoot.AddField("KM", "6.5")

	size := ""
	for _, sgf := range sgfs {
		p := New(sgf)
		root, err := p.Parse()
		if err != nil {
			// on error, just continue
			continue
		}
		eachSize := ""
		if sizes := root.GetField("SZ"); len(sizes) > 0 {
			eachSize = sizes[0]
		} else {
			// if size is not provided, assume 19
			eachSize = "19"
		}

		// if we haven't set the (assumed) same size yet, set it
		if size == "" {
			size = eachSize
		}

		// if not all the sgfs are the same size, just return the first one?
		if size != eachSize {
			return sgfs[0]
		}

		hasB := len(root.GetField("B")) > 0
		hasW := len(root.GetField("W")) > 0
		hasAB := len(root.GetField("AB")) > 0
		hasAW := len(root.GetField("AW")) > 0

		if hasB || hasW || hasAB || hasAW {
			// strip fields and save the node
			for _, key := range []string{"PB", "PW", "RE", "KM", "DT"} {
				values := root.GetField(key)
				if len(values) == 0 {
					continue
				}
				v := fmt.Sprintf("%s: %s", key, values[0])
				root.AddField("C", v)
			}
			for _, key := range []string{"RU", "SZ", "KM", "TM", "OT"} {
				root.DeleteField(key)
			}
			newRoot.down = append(newRoot.down, root)
		} else {
			// otherwise save all the children
			for _, d := range root.down {
				for _, key := range []string{"PB", "PW", "RE", "KM", "DT"} {
					if len(root.GetField(key)) == 0 {
						continue
					}
					rootValues := root.GetField(key)
					if len(rootValues) == 0 {
						continue
					}
					rootValue := root.GetField(key)[0]
					value := fmt.Sprintf("%s: %s", key, rootValue)
					d.AddField("C", value)
				}
				newRoot.down = append(newRoot.down, d)
			}
		}
	}

	newRoot.AddField("SZ", size)
	return newRoot.toSGF(true)
}
