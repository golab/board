package core

import (
	"fmt"
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
