/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core

import (
	"fmt"
	"strconv"
	"strings"
)

type Color int

const (
	NoColor Color = iota
	Black
	White
)

type FrameType int

const (
	DiffFrame = iota
	FullFrame
)

type Frame struct {
	Type     FrameType `json:"type"`
	Diff     *Diff     `json:"diff"`
	Marks    *Marks    `json:"marks"`
	Comments []string  `json:"comments"`
	Metadata *Metadata `json:"metadata"`
	TreeJSON *TreeJSON `json:"tree"`
}

type Marks struct {
	Current   *Coord   `json:"current"`
	Squares   []*Coord `json:"squares"`
	Triangles []*Coord `json:"triangles"`
	Labels    []*Label `json:"labels"`
	Pens      []*Pen   `json:"pens"`
}

type Label struct {
	Coord *Coord `json:"coord"`
	Text  string `json:"text"`
}

type Pen struct {
	X0    float64 `json:"x0"`
	Y0    float64 `json:"y0"`
	X1    float64 `json:"x1"`
	Y1    float64 `json:"y1"`
	Color string  `json:"color"`
}

type Metadata struct {
	Size   int                 `json:"size"`
	Fields map[string][]string `json:"fields"`
}

func Opposite(c Color) Color {
	if c == Black {
		return White
	}
	if c == White {
		return Black
	}
	return NoColor
}

func (c Color) String() string {
	if c == Black {
		return "B"
	}
	if c == White {
		return "W"
	}
	return "+"
}

type CoordSet map[string]*Coord

func (cs CoordSet) Has(c *Coord) bool {
	_, ok := cs[c.ToLetters()]
	return ok
}

func (cs CoordSet) Add(c *Coord) {
	cs[c.ToLetters()] = c
}

func (cs CoordSet) String() string {
	s := "["
	for k := range cs {
		s += k
		s += " "
	}
	s += "]"
	return s
}

func (cs CoordSet) List() []*Coord {
	l := []*Coord{}
	for _, c := range cs {
		l = append(l, c)
	}
	return l
}

func (cs CoordSet) IsSubsetOf(other CoordSet) bool {
	for _, v := range cs {
		if !other.Has(v) {
			return false
		}
	}
	return true
}

func (cs CoordSet) Equal(other CoordSet) bool {
	return cs.IsSubsetOf(other) && other.IsSubsetOf(cs)
}

func NewCoordSet() CoordSet {
	return CoordSet(make(map[string]*Coord))
}

type StoneSet struct {
	Coords []*Coord `json:"coords"`
	Color  `json:"color"`
}

func (s *StoneSet) Copy() *StoneSet {
	if s == nil {
		return nil
	}
	coords := []*Coord{}
	for _, c := range s.Coords {
		coords = append(coords, c.Copy())
	}
	return &StoneSet{coords, s.Color}
}

func (s *StoneSet) String() string {
	return fmt.Sprintf("%v - %v", s.Coords, s.Color)
}

func NewStoneSet(s CoordSet, c Color) *StoneSet {
	return &StoneSet{s.List(), c}
}

type Diff struct {
	Add    []*StoneSet `json:"add"`
	Remove []*StoneSet `json:"remove"`
}

func NewDiff(add, remove []*StoneSet) *Diff {
	return &Diff{
		Add:    add,
		Remove: remove,
	}
}

func (d *Diff) Copy() *Diff {
	if d == nil {
		return nil
	}
	add := []*StoneSet{}
	remove := []*StoneSet{}
	for _, a := range d.Add {
		add = append(add, a.Copy())
	}
	for _, r := range d.Remove {
		remove = append(remove, r.Copy())
	}
	return NewDiff(add, remove)
}

func (d *Diff) Invert() *Diff {
	if d == nil {
		return nil
	}
	return NewDiff(d.Remove, d.Add)
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c *Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

func (c *Coord) ToLetters() string {
	alphabet := "abcdefghijklmnopqrs"
	return string([]byte{alphabet[c.X], alphabet[c.Y]})
}

func (c *Coord) Equal(other *Coord) bool {
	if c == nil || other == nil {
		return false
	}
	return c.X == other.X && c.Y == other.Y
}

func (c *Coord) Copy() *Coord {
	if c == nil {
		return nil
	}
	return &Coord{c.X, c.Y}
}

func LettersToCoord(s string) *Coord {
	if len(s) != 2 {
		return nil
	}
	t := strings.ToLower(s)
	return &Coord{int(t[0] - 97), int(t[1] - 97)}
}

func InterfaceToCoord(ifc interface{}) (*Coord, error) {
	coords := make([]int, 0)

	// coerce the value to an array
	val, ok := ifc.([]interface{})

	if !ok {
		return nil, fmt.Errorf("error coercing to coord")
	}

	for _, v := range val {
		i := int(v.(float64))
		coords = append(coords, i)
	}
	x := coords[0]
	y := coords[1]
	return &Coord{x, y}, nil
}

type TreeJSONType int

const (
	CurrentOnly TreeJSONType = iota
	CurrentAndPreferred
	PartialNodes
	Full
)

type NodeJSON struct {
	Color Color `json:"color"`
	Down  []int `json:"down"`
	Depth int   `json:"depth"`
}

type TreeJSON struct {
	Nodes     map[int]*NodeJSON `json:"nodes"`
	Current   int               `json:"current"`
	Preferred []int             `json:"preferred"`
	Depth     int               `json:"depth"`
	Up        int               `json:"up"`
	Root      int               `json:"root"`
}

type PatternMove struct {
	Coord *Coord // nil for passes
	Color Color
}

type EventJSON struct {
	Event     string      `json:"event"`
	Value     interface{} `json:"value"`
	Color     int         `json:"color"`
	UserID    string      `json:"userid"`
	Signature string      `json:"signature"`
}

func ErrorJSON(msg string) *EventJSON {
	return &EventJSON{"error", msg, 0, "", ""}
}

func FrameJSON(frame *Frame) *EventJSON {
	return &EventJSON{"frame", frame, 0, "", ""}
}

func NopJSON() *EventJSON {
	return &EventJSON{"nop", nil, 0, "", ""}
}

func AlphanumericToCoord(s string) (*Coord, error) {
	s = strings.ToLower(s)
	if len(s) < 2 {
		return nil, fmt.Errorf("failure to parse: %s", s)
	}
	letter := s[0]
	if letter < 'a' || letter > 't' || letter == 'j' {
		return nil, fmt.Errorf("bad character: %c", letter)
	}

	num64, err := strconv.ParseInt(s[1:], 10, 64)
	if err != nil {
		return nil, err
	}
	num := int(num64)
	if num < 1 || num > 19 {
		return nil, fmt.Errorf("bad number: %d", num)
	}
	y := 18 - (num - 1)
	var x int
	if letter < 'j' {
		x = int(letter - 'a')
	} else {
		x = int(letter - 'a' - 1)
	}

	return &Coord{x, y}, nil
}
