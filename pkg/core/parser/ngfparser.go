/*
Copyright (c) 2026 Jared Nishikawa

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

type NGFParser struct {
	*BaseParser
}

func NewNGFParser(text string) *NGFParser {
	return &NGFParser{&BaseParser{[]rune(text), 0}}
}

func NGFParserFrom(p *BaseParser) *NGFParser {
	return &NGFParser{p}
}

type userdata struct {
	nick string
	rank string
}

func (p *NGFParser) parseNickRank() (*userdata, error) {
	s := p.parseLine()
	tokens := strings.Split(s, " ")
	if len(tokens) < 2 {
		return nil, fmt.Errorf("error parsing userdata")
	}
	rank := tokens[len(tokens)-1]
	nick := strings.TrimSpace(strings.Join(tokens[:len(tokens)-1], " "))
	return &userdata{nick: nick, rank: rank}, nil
}

type move struct {
	key   string
	coord string
}

func (p *NGFParser) parseMove(size int) (*move, error) {
	line := p.parseLine()
	if len(line) != 9 {
		return nil, fmt.Errorf("error parsing move")
	}
	if line[:2] != "PM" {
		return nil, fmt.Errorf("unknown key in move line")
	}
	key := fmt.Sprintf("%c", line[4])
	coordOrig := strings.ToLower(line[5:7])
	// black magic to fix coordinate orientation
	// based on extensive research (one sample)
	// wbaduk displays the board with the "origin" in the bottom right
	x := byte(size) + 'a' - (coordOrig[0] - 'a')
	y := byte(size) + 'a' - (coordOrig[1] - 'a')
	coord := fmt.Sprintf("%c%c", x, y)
	p.skipWhitespace()
	return &move{key: key, coord: coord}, nil
}

func (p *NGFParser) parseResult() string {
	lineOrig := p.parseLine()
	line := strings.ToLower(lineOrig)
	if len(line) < 5 {
		return ""
	}
	var key string
	var r string
	if strings.Contains(line[:5], "black") {
		key = "B"
	} else if strings.Contains(line[:5], "white") {
		key = "W"
	} else {
		return ""
	}
	if strings.Contains(line, "time") {
		r = "T"
	} else if strings.Contains(line, "resign") {
		r = "R"
	} else if strings.Contains(line, "points") {
		tokens := strings.Split(line, " ")
		if len(tokens) < 4 {
			r = ""
		} else {
			r = tokens[3]
		}
	}

	return key + "+" + r
}

type NGFResult struct {
	title     string
	size      int
	userWhite *userdata
	userBlack *userdata
	website   string
	handicap  int
	komi      float64
	date      string
	result    string

	moves []*move
}

func (p *NGFParser) Parse() (*NGFResult, error) {
	p.skipWhitespace()

	title := p.parseLine()

	size, err := p.parseInt()
	if err != nil {
		return nil, err
	}

	userWhite, err := p.parseNickRank()
	if err != nil {
		return nil, err
	}

	userBlack, err := p.parseNickRank()
	if err != nil {
		return nil, err
	}

	website := p.parseLine()

	handicap, err := p.parseInt()
	if err != nil {
		return nil, err
	}

	p.parseLine()

	komiInt, err := p.parseInt()
	if err != nil {
		return nil, err
	}
	komi := float64(komiInt) + 0.5

	date := p.parseLine()

	// ignore next line
	p.parseLine()

	result := p.parseResult()

	numMoves, err := p.parseInt()
	if err != nil {
		return nil, err
	}

	if numMoves < 0 {
		return nil, fmt.Errorf("numMoves should be >= 0")
	}
	if numMoves > (1 << 16) {
		return nil, fmt.Errorf("numMoves should be < 2^16")
	}
	moves := make([]*move, 0, numMoves)
	for {
		c := p.peek(0)
		if c != 'P' {
			break
		}
		move, err := p.parseMove(size)
		if err != nil {
			continue
		}
		moves = append(moves, move)
	}
	ngfResult := &NGFResult{
		title:     title,
		size:      size,
		userWhite: userWhite,
		userBlack: userBlack,
		website:   website,
		handicap:  handicap,
		komi:      komi,
		date:      date,
		result:    result,
		moves:     moves,
	}
	return ngfResult, nil
}

func (r *NGFResult) ToSGFNode() (*SGFNode, error) {
	root := &SGFNode{}
	root.AddField("SZ", strconv.Itoa(r.size))
	if r.userBlack != nil {
		root.AddField("PB", r.userBlack.nick)
		root.AddField("BR", r.userBlack.rank)
	}
	if r.userWhite != nil {
		root.AddField("PW", r.userWhite.nick)
		root.AddField("WR", r.userWhite.rank)
	}
	root.AddField("KM", strconv.FormatFloat(r.komi, 'f', -1, 64))
	root.AddField("DT", r.date)
	if r.result != "" {
		root.AddField("RE", r.result)
	}
	root.AddField("GN", r.title)
	root.AddField("PC", r.website)
	addHandicap(root, r.handicap)
	cur := root
	for _, move := range r.moves {
		node := &SGFNode{}
		node.AddField(move.key, move.coord)
		cur.down = append(cur.down, node)
		cur = node
	}
	return root, nil
}
