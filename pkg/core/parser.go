/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core

import (
	"fmt"
	"sort"
	"strings"
)

func IsWhitespace(c rune) bool {
	return c == '\n' || c == ' ' || c == '\t' || c == '\r'
}

type SGFNode struct {
	Fields map[string][]string
	Down   []*SGFNode
	Index  int
}

func (n *SGFNode) IsMove() bool {
	if _, ok := n.Fields["B"]; ok {
		return true
	}
	if _, ok := n.Fields["W"]; ok {
		return true
	}
	return false
}

func (n *SGFNode) IsPass() bool {
	if val, ok := n.Fields["B"]; ok {
		return len(val) == 1 && val[0] == ""
	}
	if val, ok := n.Fields["W"]; ok {
		return len(val) == 1 && val[0] == ""
	}
	return false
}

func (n *SGFNode) Color() Color {
	if _, ok := n.Fields["B"]; ok {
		return Black
	}
	if _, ok := n.Fields["W"]; ok {
		return White
	}
	return NoColor
}

func (n *SGFNode) Coord() *Coord {
	if val, ok := n.Fields["B"]; ok {
		if len(val) == 1 {
			return LettersToCoord(val[0])
		}
	}
	if val, ok := n.Fields["W"]; ok {
		if len(val) == 1 {
			return LettersToCoord(val[0])
		}
	}
	return nil
}

func (n *SGFNode) ToSGF(root bool) string {
	result := ""
	if root {
		result += "("
	}
	result += ";"
	fields := []string{}
	for f := range n.Fields {
		fields = append(fields, f)
	}

	sort.Strings(fields)

	for _, field := range fields {
		values := n.Fields[field]
		result += field
		for _, value := range values {
			result += "["
			result += strings.ReplaceAll(value, "]", "\\]")
			result += "]"
		}
	}

	for _, d := range n.Down {
		if len(n.Down) > 1 {
			result += "(" + d.ToSGF(false) + ")"
		} else {
			result += d.ToSGF(false)
		}
	}
	if root {
		result += ")"
	}
	return result
}

func NewSGFNode(fields map[string][]string, index int) *SGFNode {
	return &SGFNode{fields, []*SGFNode{}, index}
}

type Parser struct {
	Text  []rune
	Index int
}

func NewParser(text string) *Parser {
	return &Parser{[]rune(text), 0}
}

func (p *Parser) Parse() (*SGFNode, error) {
	p.SkipWhitespace()
	c := p.read()
	if c == '(' {
		root, err := p.ParseBranch()
		if err != nil {
			return nil, err
		}
		return Validate(root)
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
	s := ""
	for {
		t := p.read()
		if t == 0 {
			return "", fmt.Errorf("bad field")
		} else if t == ']' {
			break
		} else if t == '\\' && p.peek(0) == ']' {
			t = p.read()
		}
		s = fmt.Sprintf("%s%c", s, t)
	}
	return s, nil
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
			cur.Down = append(cur.Down, next)
			cur = next
		} else {
			break
		}
	}
	return []*SGFNode{root, cur}, nil
}

func (p *Parser) ParseNode() (*SGFNode, error) {
	fields := make(map[string][]string)
	index := 0
	for {
		p.SkipWhitespace()
		c := p.peek(0)
		if c == '(' || c == ';' || c == ')' {
			break
		}
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("bad node (expected key) %c", c)
		}
		key, err := p.ParseKey()
		if err != nil {
			return nil, err
		}
		multifield := []string{}
		p.SkipWhitespace()
		if p.read() != '[' {
			return nil, fmt.Errorf("bad node (expected field) %c", c)
		}
		field, err := p.ParseField()
		if err != nil {
			return nil, err
		}
		multifield = append(multifield, field)

		for {
			p.SkipWhitespace()
			if p.peek(0) == '[' {
				p.read()
				field, err = p.ParseField()
				if err != nil {
					return nil, err
				}
				multifield = append(multifield, field)
			} else {
				break
			}
		}

		p.SkipWhitespace()
		fields[key] = multifield
	}

	n := NewSGFNode(fields, index)
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
				current.Down = append(current.Down, node)
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
				current.Down = append(current.Down, newBranch)
			}
		} else if c == ')' {
			break
		}
	}
	return root, nil
}

func (p *Parser) read() rune {
	if p.Index >= len(p.Text) {
		return 0
	}
	result := p.Text[p.Index]
	p.Index++
	return result
}

func (p *Parser) peek(n int) rune {
	if p.Index+n >= len(p.Text) {
		return 0
	}
	return p.Text[p.Index+n]
}

func Merge(sgfs []string) string {
	if len(sgfs) == 0 {
		return ""
	} else if len(sgfs) == 1 {
		return sgfs[0]
	}

	size := ""
	fields := make(map[string][]string)
	fields["GM"] = []string{"1"}
	fields["FF"] = []string{"4"}
	fields["CA"] = []string{"UTF-8"}
	fields["PB"] = []string{"Black"}
	fields["PW"] = []string{"White"}
	fields["RU"] = []string{"Japanese"}
	fields["KM"] = []string{"6.5"}

	newRoot := NewSGFNode(fields, 0)

	for _, sgf := range sgfs {
		p := NewParser(sgf)
		root, err := p.Parse()
		if err != nil {
			// on error, just continue
			continue
		}
		eachSize := ""
		if sizes, ok := root.Fields["SZ"]; ok {
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

		_, hasB := root.Fields["B"]
		_, hasW := root.Fields["W"]
		_, hasAB := root.Fields["AB"]
		_, hasAW := root.Fields["AW"]

		if hasB || hasW || hasAB || hasAW {
			// strip fields and save the node
			for _, key := range []string{"RU", "SZ", "KM", "TM", "OT"} {
				delete(root.Fields, key)
			}
			newRoot.Down = append(newRoot.Down, root)
		} else {
			// otherwise save all the children
			for _, d := range root.Down {
				if _, ok := d.Fields["C"]; !ok {
					d.Fields["C"] = []string{}
				}
				for _, key := range []string{"PB", "PW", "RE", "KM", "DT"} {
					if len(root.Fields[key]) == 0 {
						continue
					}
					value := root.Fields[key][0]
					d.Fields["C"] = append(d.Fields["C"], fmt.Sprintf("%s: %s", key, value))
				}
				if len(d.Fields["C"]) == 0 {
					delete(d.Fields, "C")
				}
				newRoot.Down = append(newRoot.Down, d)
			}
		}
	}

	newRoot.Fields["SZ"] = []string{size}
	return newRoot.ToSGF(true)
}

func Validate(node *SGFNode) (*SGFNode, error) {
	if node == nil {
		return nil, nil
	}
	fields := make(map[string][]string)
	for key, value := range node.Fields {
		if (key == "B" || key == "W") && len(value) == 1 && value[0] == "tt" {
			fields[key] = []string{""}
		} else {
			fields[key] = value
		}
	}

	down := []*SGFNode{}
	for _, d := range node.Down {
		e, err := Validate(d)
		if err != nil {
			return nil, err
		}
		down = append(down, e)
	}

	n := &SGFNode{
		Fields: fields,
		Down:   down,
		Index:  node.Index,
	}

	return n, nil
}
