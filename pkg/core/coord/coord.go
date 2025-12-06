/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package coord

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jarednogo/board/pkg/core/color"
)

const maxBoardSize = 19

var allCoords [maxBoardSize][maxBoardSize]*Coord

func init() {
	for r := 0; r < maxBoardSize; r++ {
		for c := 0; c < maxBoardSize; c++ {
			allCoords[r][c] = &Coord{X: r, Y: c}
		}
	}
}

// CoordSet is used for quickly checking a set for existence of coords
type CoordSet map[int]*Coord

// Has: uses ToLetters as the key
func (cs CoordSet) Has(c *Coord) bool {
	_, ok := cs[c.Index()]
	return ok
}

// Add adds a coord to the set
func (cs CoordSet) Add(c *Coord) {
	cs[c.Index()] = c
}

func (cs CoordSet) AddAll(ds CoordSet) {
	for _, d := range ds {
		cs.Add(d)
	}
}

func (cs CoordSet) Remove(c *Coord) {
	delete(cs, c.Index())
}

func (cs CoordSet) RemoveAll(ds CoordSet) {
	for _, d := range ds {
		cs.Remove(d)
	}
}

func (cs CoordSet) Intersect(ds CoordSet) CoordSet {
	result := NewCoordSet()
	for k, v := range ds {
		if _, ok := cs[k]; ok {
			result.Add(v)
		}
	}
	return result
}

// String is for debugging
func (cs CoordSet) String() string {
	s := "["
	for k := range cs {
		s += indexToCoord(k).ToLetters()
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
	return CoordSet(make(map[int]*Coord))
}

// StoneSet is an array of Coords plus a Color
type StoneSet struct {
	Coords      []*Coord `json:"coords"`
	color.Color `json:"color"`
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

func (s *StoneSet) Equal(other *StoneSet) bool {
	a := NewCoordSet()
	b := NewCoordSet()
	for _, c := range s.Coords {
		a.Add(c)
	}
	for _, c := range other.Coords {
		b.Add(c)
	}
	return a.Equal(b) && s.Color == other.Color
}

// String is for debugging
func (s *StoneSet) String() string {
	return fmt.Sprintf("%v - %v", s.Coords, s.Color)
}

// NewStoneSet takes a CoordSet and a Color and turns it into a StoneSet
func NewStoneSet(s CoordSet, c color.Color) *StoneSet {
	return &StoneSet{s.List(), c}
}

// Coord is just an (x,y) coordinate pair
type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewCoord(x, y int) *Coord {
	if x < 0 || x >= len(allCoords) {
		return nil
	}
	if y < 0 || y >= len(allCoords[x]) {
		return nil
	}
	return allCoords[x][y]
}

func indexToCoord(i int) *Coord {
	return NewCoord(i/maxBoardSize, i%maxBoardSize)
}

func (c *Coord) Index() int {
	return c.X*maxBoardSize + c.Y
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
	if c == nil && other == nil {
		return true
	}
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
	return NewCoord(c.X, c.Y)
}

// LettersToCoord takes a pair of letters and turns it into a Coord
func LettersToCoord(s string) *Coord {
	if len(s) != 2 {
		return nil
	}
	t := strings.ToLower(s)
	return NewCoord(int(t[0]-97), int(t[1]-97))
}

// InterfaceToCoord essentially coerces the interface into a [2]int
// then turns that into a Coord
func InterfaceToCoord(ifc any) (*Coord, error) {
	coords := make([]int, 0)

	// coerce the value to an array
	val, ok := ifc.([]any)

	if !ok {
		return nil, fmt.Errorf("error coercing to coord")
	}

	for _, v := range val {
		i := int(v.(float64))
		coords = append(coords, i)
	}
	x := coords[0]
	y := coords[1]
	return NewCoord(x, y), nil
}

// Stone is a Coord plus a Color
// StoneSet is essentially an extension of Stone
type Stone struct {
	Coord *Coord // nil for passes
	Color color.Color
}

func NewStone(x, y int, c color.Color) *Stone {
	coord := NewCoord(x, y)
	return &Stone{
		Coord: coord,
		Color: c,
	}
}

// AlphanumericToCoord converts a string like "c17" to a Coord
func AlphanumericToCoord(s string, size int) (*Coord, error) {
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
	if num < 1 || num > size {
		return nil, fmt.Errorf("bad number: %d", num)
	}
	y := size - num
	var x int
	if letter < 'i' {
		x = int(letter - 'a')
	} else {
		x = int(letter - 'a' - 1)
	}

	return NewCoord(x, y), nil
}
