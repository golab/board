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

type userdata struct {
	nick string
	rank string
}

func (p *NGFParser) parseNickRank() (*userdata, error) {
	s, err := p.parseLine()
	if err != nil {
		return nil, err
	}
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

func (p *NGFParser) parseMove() (*move, error) {
	line, err := p.parseLine()
	if err != nil {
		return nil, err
	}
	if len(line) != 9 {
		return nil, fmt.Errorf("error parsing move")
	}
	if line[:2] != "PM" {
		return nil, fmt.Errorf("unknown key in move line")
	}
	key := fmt.Sprintf("%c", line[4])
	coordOrig := strings.ToLower(line[5:7])
	x := coordOrig[0] - 1
	y := coordOrig[1] - 1
	coord := fmt.Sprintf("%c%c", x, y)
	p.skipWhitespace()
	return &move{key: key, coord: coord}, nil
}

func (p *NGFParser) parseResult() (string, error) {
	lineOrig, err := p.parseLine()
	if err != nil {
		return "", err
	}
	line := strings.ToLower(lineOrig)
	var key string
	var r string
	if strings.Contains(line[:5], "black") {
		key = "B"
	} else if strings.Contains(line[:5], "white") {
		key = "W"
	} else {
		return "", nil
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

	return key + "+" + r, nil
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

	title, err := p.parseLine()
	if err != nil {
		return nil, err
	}

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

	website, err := p.parseLine()
	if err != nil {
		return nil, err
	}

	handicap, err := p.parseInt()
	if err != nil {
		return nil, err
	}

	_, err = p.parseLine()
	if err != nil {
		return nil, err
	}

	komiInt, err := p.parseInt()
	if err != nil {
		return nil, err
	}
	komi := float64(komiInt) + 0.5

	date, err := p.parseLine()
	if err != nil {
		return nil, err
	}

	_, err = p.parseLine()
	if err != nil {
		return nil, err
	}

	result, err := p.parseResult()
	if err != nil {
		return nil, err
	}

	numMoves, err := p.parseInt()
	if err != nil {
		return nil, err
	}

	moves := make([]*move, numMoves)
	index := 0
	for {
		c := p.peek(0)
		if c != 'P' {
			break
		}
		move, err := p.parseMove()
		if err != nil {
			continue
		}
		if index < len(moves) {
			moves[index] = move
			index++
		} else {
			moves = append(moves, move)
		}
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
	root.AddField("RE", r.result)
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
