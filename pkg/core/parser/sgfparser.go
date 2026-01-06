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
)

type SGFParser struct {
	*BaseParser
}

func NewSGFParser(text string) *SGFParser {
	return &SGFParser{&BaseParser{[]rune(text), 0}}
}

func (p *SGFParser) Parse() (*SGFNode, error) {
	// now parse sgf
	_, err := p.parseUntil('(')
	if err != nil {
		return nil, err
	}
	// the next character is guaranteed to be '('
	root, err := p.parseBranch()
	if err != nil {
		return nil, err
	}
	return root, nil
}

func (p *SGFParser) parseKey() (string, error) {
	sb := strings.Builder{}
	for {
		c := p.peek(0)
		if isLower(c) {
			r := p.read()
			r -= 32
			sb.WriteRune(r)
		} else if isUpper(c) {
			sb.WriteRune(p.read())
		} else if c == '[' {
			break
		} else if isWhitespace(c) {
			p.skipWhitespace()
		} else {
			return "", fmt.Errorf("bad key %c", c)
		}
	}
	return sb.String(), nil
}

func (p *SGFParser) parseField() (string, error) {
	// read '['
	err := p.requireRune('[')
	if err != nil {
		return "", err
	}
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

type parseNodesResult struct {
	root *SGFNode
	cur  *SGFNode
}

func (p *SGFParser) parseOneOrMoreNodes() (*parseNodesResult, error) {
	n, err := p.parseNode()
	if err != nil {
		return nil, err
	}
	root := n
	cur := root
	for {
		c := p.peek(0)
		if c == ';' {
			next, err := p.parseNode()
			if err != nil {
				return nil, err
			}
			cur.down = append(cur.down, next)
			cur = next
		} else {
			break
		}
	}
	return &parseNodesResult{root, cur}, nil
}

type property struct {
	key    string
	values []string
}

func (p *SGFParser) parseProperty() (*property, error) {
	prop := &property{}
	c := p.peek(0)
	if !isLower(c) && !isUpper(c) {
		return nil, fmt.Errorf("bad property (expected key) %c", c)
	}
	key, err := p.parseKey()
	if err != nil {
		return nil, err
	}

	prop.key = key

	p.skipWhitespace()

	flds, err := p.parseOneOrMoreFields(key)
	if err != nil {
		return nil, err
	}

	prop.values = flds

	return prop, nil
}

func (p *SGFParser) parseOneOrMoreFields(key string) ([]string, error) {
	flds := []string{}
	// require parse first field
	field, err := p.parseOneField(key)
	if err != nil {
		return nil, err
	}

	flds = append(flds, field)

	// potentially parse more fields
	for {
		p.skipWhitespace()
		if p.peek(0) == '[' {
			field, err := p.parseOneField(key)
			if err != nil {
				return nil, err
			}
			flds = append(flds, field)
		} else {
			break
		}
	}
	return flds, nil
}

func (p *SGFParser) parseOneField(key string) (string, error) {
	field, err := p.parseField()
	if err != nil {
		return "", err
	}

	if (key == "B" || key == "W") && field == "tt" {
		field = ""
	}
	return field, nil
}

func (p *SGFParser) parseNode() (*SGFNode, error) {
	// read ';'
	err := p.requireRune(';')
	if err != nil {
		return nil, err
	}
	n := &SGFNode{}
	for {
		p.skipWhitespace()
		c := p.peek(0)
		if c == '(' || c == ';' || c == ')' {
			break
		}

		prop, err := p.parseProperty()
		if err != nil {
			return nil, err
		}

		n.SetField(prop.key, prop.values)

		p.skipWhitespace()
	}

	return n, nil
}

func (p *SGFParser) parseBranch() (*SGFNode, error) {
	// read '('
	err := p.requireRune('(')
	if err != nil {
		return nil, err
	}
	var root *SGFNode
	var current *SGFNode
	for {
		p.skipWhitespace()
		c := p.peek(0)
		if c == 0 {
			return nil, fmt.Errorf("unfinished branch, expected ')'")
		} else if c == ';' {
			result, err := p.parseOneOrMoreNodes()
			if err != nil {
				return nil, err
			}
			root = result.root
			current = result.cur
		} else if c == '(' {
			newBranch, err := p.parseBranch()
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
			// consume ')'
			p.read()
			break
		} else {
			return nil, fmt.Errorf("improperly formatted branch %c", c)
		}
	}
	return root, nil
}
