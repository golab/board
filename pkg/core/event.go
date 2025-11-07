/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// the core package provides basic functionality to all the major components of the code
package core

// Event is the basic struct for sending and receiving messages over
// the websockets
type Event struct {
	Type   string      `json:"event"`
	Value  interface{} `json:"value"`
	UserID string      `json:"userid"`
}

func EmptyEvent() *Event {
	return &Event{}
}

// ErrorEvent is a special case of an Event
func ErrorEvent(msg string) *Event {
	return &Event{
		Type:   "error",
		Value:  msg,
		UserID: ""}
}

// FrameEvent is a special case of an Event
func FrameEvent(frame *Frame) *Event {
	return &Event{
		Type:   "frame",
		Value:  frame,
		UserID: ""}
}

// NopEvent signals to the server to do nothing
// (in particular, don't send to clients)
func NopEvent() *Event {
	return &Event{
		Type:   "nop",
		Value:  nil,
		UserID: ""}
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

// FrameType can be either DiffFrame or FullFrame
type FrameType int

const (
	DiffFrame = iota
	FullFrame
)

// Frame provides the data for when the board needs to be updated (not the explorer)
type Frame struct {
	Type      FrameType `json:"type"`
	Diff      *Diff     `json:"diff"`
	Marks     *Marks    `json:"marks"`
	Comments  []string  `json:"comments"`
	Metadata  *Metadata `json:"metadata"`
	TreeJSON  *TreeJSON `json:"tree"`
	BlackCaps int       `json:"black_caps"`
	WhiteCaps int       `json:"white_caps"`
	BlackArea []*Coord  `json:"black_area"`
	WhiteArea []*Coord  `json:"white_area"`
	Dame      []*Coord  `json:"dame"`
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
