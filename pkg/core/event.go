/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// the core package provides basic functionality to all the major components of the code
package core

// EventJSON is the basic struct for sending and receiving messages over
// the websockets
type EventJSON struct {
	Event  string      `json:"event"`
	Value  interface{} `json:"value"`
	UserID string      `json:"userid"`
}

func EmptyEvent() *EventJSON {
	return &EventJSON{}
}

// ErrorEvent is a special case of an EventJSON
func ErrorEvent(msg string) *EventJSON {
	return &EventJSON{"error", msg, ""}
}

// FrameEvent is a special case of an EventJSON
func FrameEvent(frame *Frame) *EventJSON {
	return &EventJSON{"frame", frame, ""}
}

// NopEvent signals to the server to do nothing
// (in particular, don't send to clients)
func NopEvent() *EventJSON {
	return &EventJSON{"nop", nil, ""}
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
