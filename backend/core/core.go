/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// the core package provides basic functionality to all the major components of the code
package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Color is one of NoColor, Black, or White
type Color int

const (
	NoColor Color = iota
	Black
	White
)

// FrameType can be either DiffFrame or FullFrame
type FrameType int

const (
	DiffFrame = iota
	FullFrame
)

// Frame provides the data for when the board needs to be updated (not the explorer)
type Frame struct {
	Type     FrameType `json:"type"`
	Diff     *Diff     `json:"diff"`
	Marks    *Marks    `json:"marks"`
	Comments []string  `json:"comments"`
	Metadata *Metadata `json:"metadata"`
	TreeJSON *TreeJSON `json:"tree"`
}

// Marks provides data for any marks on the board
type Marks struct {
	Current   *Coord   `json:"current"`
	Squares   []*Coord `json:"squares"`
	Triangles []*Coord `json:"triangles"`
	Labels    []*Label `json:"labels"`
	Pens      []*Pen   `json:"pens"`
}

// Label can be any text, but typically single digits or letters
type Label struct {
	Coord *Coord `json:"coord"`
	Text  string `json:"text"`
}

// Pen contains a start and end coordinate plus a color
type Pen struct {
	X0    float64 `json:"x0"`
	Y0    float64 `json:"y0"`
	X1    float64 `json:"x1"`
	Y1    float64 `json:"y1"`
	Color string  `json:"color"`
}

// Metadata provides the size of the board plus any fields (usually from the root node)
type Metadata struct {
	Size   int                 `json:"size"`
	Fields map[string][]string `json:"fields"`
}

// Opposite: Black -> White, White -> Black, NoColor -> NoColor
func Opposite(c Color) Color {
	if c == Black {
		return White
	}
	if c == White {
		return Black
	}
	return NoColor
}

// String is just for debugging purposes
func (c Color) String() string {
	if c == Black {
		return "B"
	}
	if c == White {
		return "W"
	}
	return "+"
}

// CoordSet is used for quickly checking a set for existence of coords
type CoordSet map[string]*Coord

// Has: uses ToLetters as the key
func (cs CoordSet) Has(c *Coord) bool {
	_, ok := cs[c.ToLetters()]
	return ok
}

// Add adds a coord to the set
func (cs CoordSet) Add(c *Coord) {
	cs[c.ToLetters()] = c
}

// String is for debugging
func (cs CoordSet) String() string {
	s := "["
	for k := range cs {
		s += k
		s += " "
	}
	s += "]"
	return s
}

// List converts the map to an array
func (cs CoordSet) List() []*Coord {
	l := []*Coord{}
	for _, c := range cs {
		l = append(l, c)
	}
	return l
}

// IsSubsetOf checks for set inclusion
func (cs CoordSet) IsSubsetOf(other CoordSet) bool {
	for _, v := range cs {
		if !other.Has(v) {
			return false
		}
	}
	return true
}

// Equal does IsSubsetOf twice
func (cs CoordSet) Equal(other CoordSet) bool {
	return cs.IsSubsetOf(other) && other.IsSubsetOf(cs)
}

// NewCoordSet makes an empty CoordSet
func NewCoordSet() CoordSet {
	return CoordSet(make(map[string]*Coord))
}

// StoneSet is an array of Coords plus a Color
type StoneSet struct {
	Coords []*Coord `json:"coords"`
	Color  `json:"color"`
}

// Copy copies the StoneSet
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

// String is for debugging
func (s *StoneSet) String() string {
	return fmt.Sprintf("%v - %v", s.Coords, s.Color)
}

// NewStoneSet takes a CoordSet and a Color and turns it into a StoneSet
func NewStoneSet(s CoordSet, c Color) *StoneSet {
	return &StoneSet{s.List(), c}
}

// Diff contains two StoneSets (Add and Remove) and is a key component of a Frame
type Diff struct {
	Add    []*StoneSet `json:"add"`
	Remove []*StoneSet `json:"remove"`
}

// NewDiff makes a Diff based on two StoneSets
func NewDiff(add, remove []*StoneSet) *Diff {
	return &Diff{
		Add:    add,
		Remove: remove,
	}
}

// Copy makes a copy of the Diff
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

// Invert simply exchanges Add and Remove
func (d *Diff) Invert() *Diff {
	if d == nil {
		return nil
	}
	return NewDiff(d.Remove, d.Add)
}

// Coord is just an (x,y) coordinate pair
type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// String is for debugging
func (c *Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

// ToLetters is sgf-specific (notice the inclusion of 'i')
func (c *Coord) ToLetters() string {
	alphabet := "abcdefghijklmnopqrs"
	return string([]byte{alphabet[c.X], alphabet[c.Y]})
}

// Equal compares two Coords
func (c *Coord) Equal(other *Coord) bool {
	if c == nil || other == nil {
		return false
	}
	return c.X == other.X && c.Y == other.Y
}

// Copy copies the Coord
func (c *Coord) Copy() *Coord {
	if c == nil {
		return nil
	}
	return &Coord{c.X, c.Y}
}

// LettersToCoord takes a pair of letters and turns it into a Coord
func LettersToCoord(s string) *Coord {
	if len(s) != 2 {
		return nil
	}
	t := strings.ToLower(s)
	return &Coord{int(t[0] - 97), int(t[1] - 97)}
}

// InterfaceToCoord essentially coerces the interface into a [2]int
// then turns that into a Coord
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

// TreeJSONType defines some options for how much data to send in a TreeJSON
type TreeJSONType int

const (
	CurrentOnly TreeJSONType = iota
	CurrentAndPreferred
	PartialNodes
	Full
)

// NodeJSON is a key component of TreeJSON
type NodeJSON struct {
	Color Color `json:"color"`
	Down  []int `json:"down"`
	Depth int   `json:"depth"`
}

// TreeJSON is the basic struct to encode information about the explorer
// this makes up one part of a Frame
type TreeJSON struct {
	Nodes     map[int]*NodeJSON `json:"nodes"`
	Current   int               `json:"current"`
	Preferred []int             `json:"preferred"`
	Depth     int               `json:"depth"`
	Up        int               `json:"up"`
	Root      int               `json:"root"`
}

// PatternMove is a Coord plus a Color
// StoneSet is essentially an extension of PatternMove
// TODO: consider renaming PatternMove
type PatternMove struct {
	Coord *Coord // nil for passes
	Color Color
}

// EventJSON is the basic struct for sending and receiving messages over
// the websockets
// TODO: see if we can remove Color
type EventJSON struct {
	Event  string      `json:"event"`
	Value  interface{} `json:"value"`
	Color  int         `json:"color"`
	UserID string      `json:"userid"`
}

// ErrorJSON is a special case of an EventJSON
func ErrorJSON(msg string) *EventJSON {
	return &EventJSON{"error", msg, 0, ""}
}

// FrameJSON is a special case of an EventJSON
func FrameJSON(frame *Frame) *EventJSON {
	return &EventJSON{"frame", frame, 0, ""}
}

// NopJSON signals to the server to do nothing
// (in particular, don't send to clients)
func NopJSON() *EventJSON {
	return &EventJSON{"nop", nil, 0, ""}
}

// AlphanumericToCoord converts a string like "c17" to a Coord
func AlphanumericToCoord(s string) (*Coord, error) {
	s = strings.ToLower(s)
	if len(s) < 2 {
		return nil, fmt.Errorf("failure to parse: %s", s)
	}
	letter := s[0]
	if letter < 'a' || letter > 't' || letter == 'i' {
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
	if letter < 'i' {
		x = int(letter - 'a')
	} else {
		x = int(letter - 'a' - 1)
	}

	return &Coord{x, y}, nil
}

// MyURL gets the server url from an env var
// this will be different on test vs main
func MyURL() string {
	s := os.Getenv("MYURL")
	if s == "" {
		return "http://localhost:8080"
	}
	return s
}

// Sanitize ensures strings only contain letters and numbers
func Sanitize(s string) string {
	ok := []rune{}
	for _, c := range s {
		if (c >= '0' && c <= '9') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') {
			ok = append(ok, c)
		}
	}
	return string(ok)
}

// UUID4 makes and sanitizes a new uuid
func UUID4() string {
	r, _ := uuid.NewRandom()
	s := r.String()
	// remove hyphens
	return Sanitize(s)
}
