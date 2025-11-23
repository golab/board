/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core

import (
	"fmt"
	"strings"
)

func IsWhitespace(c rune) bool {
	return c == '\n' || c == ' ' || c == '\t' || c == '\r'
}

type SGFNode struct {
	Fields
	down []*SGFNode
}

func (n *SGFNode) GetChild(i int) *SGFNode {
	if i >= 0 && i < len(n.down) {
		return n.down[i]
	}
	return nil
}

func (n *SGFNode) NumChildren() int {
	return len(n.down)
}

func (n *SGFNode) IsMove() bool {
	bvalues := n.GetField("B")
	wvalues := n.GetField("W")
	return len(bvalues) > 0 || len(wvalues) > 0
}

func (n *SGFNode) IsPass() bool {
	bvalues := n.GetField("B")
	wvalues := n.GetField("W")
	return (len(bvalues) == 1 && bvalues[0] == "") ||
		(len(wvalues) == 1 && wvalues[0] == "")
}

func (n *SGFNode) Color() Color {
	bvalues := n.GetField("B")
	wvalues := n.GetField("W")
	if len(bvalues) > 0 {
		return Black
	}
	if len(wvalues) > 0 {
		return White
	}
	return NoColor
}

func (n *SGFNode) Coord() *Coord {
	bvalues := n.GetField("B")
	wvalues := n.GetField("W")
	if len(bvalues) == 1 {
		return LettersToCoord(bvalues[0])
	}
	if len(wvalues) == 1 {
		return LettersToCoord(wvalues[0])
	}
	return nil
}

func (n *SGFNode) ToSGF(root bool) string {
	sb := strings.Builder{}
	if root {
		sb.WriteByte('(')
	}
	sb.WriteByte(';')

	n.SortFields()

	for _, field := range n.fields {
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
			sb.WriteString(d.ToSGF(false))
			sb.WriteByte(')')
		} else {
			sb.WriteString(d.ToSGF(false))
		}
	}
	if root {
		sb.WriteByte(')')
	}
	return sb.String()
}

type Parser struct {
	text  []rune
	index int
}

func NewParser(text string) *Parser {
	return &Parser{[]rune(text), 0}
}

func (p *Parser) Parse() (*SGFNode, error) {
	p.SkipWhitespace()
	p.SkipIfNot('(')
	c := p.read()
	if c == '(' {
		root, err := p.ParseBranch()
		if err != nil {
			return nil, err
		}
		return root, nil
	}
	return nil, fmt.Errorf("unexpected %c", c)
}

func (p *Parser) SkipWhitespace() {
	for {
		if IsWhitespace(p.peek(0)) {
			p.read()
		} else {
			break
		}
	}
}

func (p *Parser) SkipIfNot(r rune) {
	for {
		c := p.peek(0)
		if c == rune(0) {
			return
		} else if c != r {
			p.read()
		} else {
			break
		}
	}
}

func (p *Parser) ParseKey() (string, error) {
	s := ""
	for {
		c := p.peek(0)
		if c == 0 {
			return "", fmt.Errorf("bad key")
		} else if c >= 'a' && c <= 'z' {
			r := p.read()
			r -= 32
			s += string([]rune{r})
		} else if c < 'A' || c > 'Z' {
			break
		} else {
			s += string([]rune{p.read()})
		}
	}
	return s, nil
}

func (p *Parser) ParseField() (string, error) {
	sb := strings.Builder{}
	for {
		t := p.read()
		if t == 0 {
			return "", fmt.Errorf("bad field")
		} else if t == ']' {
			break
		} else if t == '\\' && p.peek(0) == ']' {
			t = p.read()
		}
		sb.WriteRune(t)
	}
	return sb.String(), nil
}

func (p *Parser) ParseNodes() ([]*SGFNode, error) {
	n, err := p.ParseNode()
	if err != nil {
		return nil, err
	}
	root := n
	cur := root
	for {
		c := p.peek(0)
		if c == ';' {
			p.read()
			next, err := p.ParseNode()
			if err != nil {
				return nil, err
			}
			cur.down = append(cur.down, next)
			cur = next
		} else {
			break
		}
	}
	return []*SGFNode{root, cur}, nil
}

type property struct {
	key    string
	values []string
}

func (p *Parser) ParseProperty() (*property, error) {
	prop := &property{}
	c := p.peek(0)
	if c < 'A' || c > 'Z' {
		return nil, fmt.Errorf("bad property (expected key) %c", c)
	}
	key, err := p.ParseKey()
	if err != nil {
		return nil, err
	}

	prop.key = key

	p.SkipWhitespace()
	if p.read() != '[' {
		return nil, fmt.Errorf("bad property (expected field) %c", c)
	}

	fields, err := p.ParseOneOrMoreFields(key)
	if err != nil {
		return nil, err
	}

	prop.values = fields

	return prop, nil
}

func (p *Parser) ParseOneOrMoreFields(key string) ([]string, error) {
	fields := []string{}
	// require parse first field
	field, err := p.ParseOneField(key)
	if err != nil {
		return nil, err
	}

	fields = append(fields, field)

	// potentially parse more fields
	for {
		p.SkipWhitespace()
		if p.peek(0) == '[' {
			p.read()
			field, err := p.ParseOneField(key)
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
		} else {
			break
		}
	}
	return fields, nil
}

func (p *Parser) ParseOneField(key string) (string, error) {
	field, err := p.ParseField()
	if err != nil {
		return "", err
	}

	if (key == "B" || key == "W") && field == "tt" {
		field = ""
	}
	return field, nil
}

func (p *Parser) ParseNode() (*SGFNode, error) {
	n := &SGFNode{}
	for {
		p.SkipWhitespace()
		c := p.peek(0)
		if c == '(' || c == ';' || c == ')' {
			break
		}

		prop, err := p.ParseProperty()
		if err != nil {
			return nil, err
		}

		n.SetField(prop.key, prop.values)

		p.SkipWhitespace()
	}

	return n, nil
}

func (p *Parser) ParseBranch() (*SGFNode, error) {
	var root *SGFNode
	var current *SGFNode
	for {
		c := p.read()
		if c == 0 {
			return nil, fmt.Errorf("unfinished branch, expected ')'")
		} else if c == ';' {
			nodes, err := p.ParseNodes()
			if err != nil {
				return nil, err
			}
			node := nodes[0]
			cur := nodes[1]
			if root == nil {
				root = node
				current = cur
			} else {
				current.down = append(current.down, node)
				current = cur
			}
		} else if c == '(' {
			newBranch, err := p.ParseBranch()
			if err != nil {
				return nil, err
			}

			if root == nil {
				root = newBranch
				current = newBranch
			} else {
				current.down = append(current.down, newBranch)
			}
		} else if c == ')' {
			break
		}
	}
	return root, nil
}

func (p *Parser) read() rune {
	if p.index >= len(p.text) {
		return 0
	}
	result := p.text[p.index]
	p.index++
	return result
}

func (p *Parser) peek(n int) rune {
	if p.index+n >= len(p.text) {
		return 0
	}
	return p.text[p.index+n]
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
		p := NewParser(sgf)
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
	return newRoot.ToSGF(true)
}
