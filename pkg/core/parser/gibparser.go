/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type gameInfoParser struct {
	*BaseParser
}

func NewGameInfoParser(text string) *gameInfoParser {
	return &gameInfoParser{&BaseParser{[]rune(text), 0}}
}

func (p *gameInfoParser) parse() ([]*keyValue, error) {
	props := []*keyValue{}
	p.skipWhitespace()
	first, err := p.parseProperty()
	props = append(props, first)
	if err != nil {
		return nil, err
	}
	for {
		c := p.peek(0)
		if c == 0 {
			break
		}
		err := p.requireRune(',')
		if err != nil {
			return nil, err
		}
		prop, err := p.parseProperty()
		if err != nil {
			return nil, err
		}
		props = append(props, prop)
	}
	return props, nil
}

func (p *gameInfoParser) parseProperty() (*keyValue, error) {
	p.skipWhitespace()
	key, err := p.parseKey()
	if err != nil {
		return nil, err
	}
	err = p.requireRune(':')
	if err != nil {
		return nil, err
	}
	p.skipWhitespace()
	value := p.parseValue()
	return &keyValue{key, value}, nil
}

func (p *gameInfoParser) parseKey() (string, error) {
	return p.parseUntil(':')
}

func (p *gameInfoParser) parseValue() string {
	sb := strings.Builder{}
	for {
		c := p.peek(0)
		if c == ',' || c == 0 {
			break
		}
		sb.WriteRune(p.read())
	}
	return sb.String()
}

type GIBResult struct {
	header []*keyValue
	game   []string
}

var alphabet = "abcdefghijklmnopqrs"

func addHandicap(node *SGFNode, n int) {
	if n < 2 {
		return
	}
	node.AddField("AB", "pd")
	node.AddField("AB", "dp")
	if n > 2 {
		node.AddField("AB", "pp")
	}
	if n > 3 && n%2 == 1 {
		node.AddField("AB", "jj")
	}
	if n > 3 {
		node.AddField("AB", "dd")
	}
	if n > 5 {
		node.AddField("AB", "dj")
		node.AddField("AB", "pj")
	}
	if n > 7 {
		node.AddField("AB", "jd")
		node.AddField("AB", "jp")
	}
}

const (
	GLRTBlackByCounting = 0
	GLRTWhiteByCounting = 1
	GLRTBlackByResign   = 3
	GLRTWhiteByResign   = 4
	GLRTBlackByTime     = 7
	GLRTWhiteByTime     = 8
)

func (g *GIBResult) ToSGFNode() (*SGFNode, error) {
	root := &SGFNode{}
	for _, h := range g.header {
		switch h.key {
		case "GAMEINFOMAIN":
			props, err := NewGameInfoParser(h.value).parse()
			if err != nil {
				return nil, err
			}

			var glrt int
			for _, prop := range props {
				switch prop.key {
				case "GRLT":
					n, err := strconv.Atoi(prop.value)
					if err != nil {
						continue
					}
					switch n {
					case GLRTBlackByCounting:
						glrt = n
					case GLRTWhiteByCounting:
						glrt = n
					case GLRTBlackByResign:
						root.AddField("RE", "B+R")
					case GLRTWhiteByResign:
						root.AddField("RE", "W+R")
					case GLRTBlackByTime:
						root.AddField("RE", "B+T")
					case GLRTWhiteByTime:
						root.AddField("RE", "W+T")
					}
				case "ZIPSU":
					zipsu, err := strconv.Atoi(prop.value)
					if err != nil {
						continue
					}
					if zipsu == 0 {
						continue
					}
					f := float64(zipsu) / 10
					switch glrt {
					case GLRTBlackByCounting:
						root.AddField("RE", fmt.Sprintf("B+%f", f))
					case GLRTWhiteByCounting:
						root.AddField("RE", fmt.Sprintf("W+%f", f))
					}
				case "LINE":
					root.AddField("SZ", prop.value)
				case "GONGJE":
					n, err := strconv.Atoi(prop.value)
					if err != nil || n == 0 {
						continue
					}
					root.AddField("KM", fmt.Sprintf("%f", float64(n)/10))
				case "DUM":
					n, err := strconv.Atoi(prop.value)
					if err != nil || n == 0 {
						continue
					}
					root.AddField("KM", fmt.Sprintf("-%d", n))
				}
			}
		case "WUSERINFO", "BUSERINFO":
			props, err := NewGameInfoParser(h.value).parse()
			if err != nil {
				return nil, err
			}
			for _, prop := range props {
				switch prop.key {
				case "BNICK":
					root.AddField("PB", prop.value)
				case "WNICK":
					root.AddField("PW", prop.value)
				}
			}

		}
	}
	cur := root
	for _, line := range g.game {
		tokens := strings.Split(strings.TrimSpace(line), " ")
		if len(tokens) < 1 {
			continue
		}
		switch tokens[0] {
		case "INI":
			// handicap stones
			if len(tokens) < 4 {
				continue
			}
			num, err := strconv.Atoi(tokens[3])
			if err != nil {
				return nil, err
			}
			addHandicap(root, num)

		case "SKI":
			// pass
			key := ""
			if cur.HasField("B") {
				key = "W"
			} else if cur.HasField("W") {
				key = "B"
			} else {
				continue
			}
			node := &SGFNode{}
			node.AddField(key, "")
			cur.down = append(cur.down, node)
			cur = node
		case "STO":
			// stone
			if len(tokens) != 6 {
				continue
			}

			col, err := strconv.Atoi(tokens[3])
			if err != nil {
				return nil, err
			}
			x, err := strconv.Atoi(tokens[4])
			if err != nil {
				return nil, err
			}
			y, err := strconv.Atoi(tokens[5])
			if err != nil {
				return nil, err
			}
			node := &SGFNode{}
			key := "B"
			if col == 2 {
				key = "W"
			}
			value := string([]byte{alphabet[x], alphabet[y]})
			node.AddField(key, value)
			cur.down = append(cur.down, node)
			cur = node
		}
	}
	return root, nil
}

type GIBParser struct {
	*BaseParser
}

func NewGIBParser(text string) *GIBParser {
	return &GIBParser{&BaseParser{[]rune(text), 0}}
}

func (p *GIBParser) Parse() (*GIBResult, error) {
	header, err := p.parseHeader()
	if err != nil {
		return nil, err
	}
	game, err := p.parseGame()
	if err != nil {
		return nil, err
	}
	return &GIBResult{
		header: header,
		game:   game,
	}, nil
}

type keyValue struct {
	key   string
	value string
}

func (p *GIBParser) parseHeader() ([]*keyValue, error) {
	p.skipWhitespace()
	err := p.require("\\HS")
	if err != nil {
		return nil, err
	}
	properties := []*keyValue{}
	for {
		p.skipWhitespace()
		if p.peek(0) == '\\' && p.peek(1) == '[' {
			prop, err := p.parseProperty()
			if err != nil {
				return nil, err
			}
			properties = append(properties, prop)
		} else {
			break
		}
	}
	err = p.require("\\HE")
	if err != nil {
		return nil, err
	}
	return properties, nil
}

func (p *GIBParser) parseProperty() (*keyValue, error) {
	err := p.require("\\[")
	if err != nil {
		return nil, err
	}
	p.skipWhitespace()
	key, err := p.parseKey()
	if err != nil {
		return nil, err
	}
	err = p.require("=")
	if err != nil {
		return nil, err
	}
	p.skipWhitespace()
	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	err = p.require("\\]")
	if err != nil {
		return nil, err
	}
	return &keyValue{
		key:   key,
		value: value,
	}, nil
}

func (p *GIBParser) parseKey() (string, error) {
	return p.parseUntil('=')
}

func (p *GIBParser) parseValue() (string, error) {
	sb := strings.Builder{}
	for {
		c := p.peek(0)
		if c == 0 {
			return "", fmt.Errorf("error parsing value, encountered null")
		}
		if c == '\\' && p.peek(1) == ']' {
			break
		}
		sb.WriteRune(p.read())
	}
	return sb.String(), nil
}

func (p *GIBParser) parseGame() ([]string, error) {
	p.skipWhitespace()
	err := p.require("\\GS")
	if err != nil {
		return nil, err
	}

	lines := []string{}
	for {
		p.skipWhitespace()
		if p.peek(0) == '\\' {
			break
		}
		line, err := p.parseGameLine()
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	err = p.require("\\GE")
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func (p *GIBParser) parseGameLine() (string, error) {
	sb := strings.Builder{}
	for {
		c := p.read()
		if c == 0 {
			return "", fmt.Errorf("error parsing game line, encountered null")
		}
		if c == '\n' {
			break
		}
		sb.WriteRune(c)
	}
	return sb.String(), nil
}
