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

type BaseParser struct {
	text  []rune
	index int
}

func (p *BaseParser) parseUntil(r rune) (string, error) {
	sb := strings.Builder{}
	for {
		c := p.peek(0)
		if c == rune(0) {
			return "", fmt.Errorf("no %c found", r)
		} else if c != r {
			sb.WriteRune(c)
			p.read()
		} else {
			break
		}
	}
	return sb.String(), nil
}

func (p *BaseParser) require(s string) error {
	sb := strings.Builder{}
	for i := 0; i < len(s); i++ {
		sb.WriteRune(p.read())
	}
	t := sb.String()
	if t != s {
		return fmt.Errorf("error parsing: required %s, got %s", s, t)
	}
	return nil
}

func (p *BaseParser) requireRune(r rune) error {
	c := p.read()
	if c != r {
		return fmt.Errorf("expected %c, got %c", r, c)
	}
	return nil
}

func (p *BaseParser) read() rune {
	if p.index >= len(p.text) {
		return 0
	}
	result := p.text[p.index]
	p.index++
	return result
}

func (p *BaseParser) peek(n int) rune {
	if p.index+n >= len(p.text) {
		return 0
	}
	return p.text[p.index+n]
}

func (p *BaseParser) skipWhitespace() {
	for {
		if isWhitespace(p.peek(0)) {
			p.read()
		} else {
			break
		}
	}
}
