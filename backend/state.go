/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

const Letters = "ABCDEFGHIJKLNMOPQRSTUVWXYZ"

type Settings struct {
	Buffer   int64
	Size     int
	Password string
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

// as a rule, anything that would need to get sent to new connections
// should be stored here and not in the Room struct
type State struct {
	Root        *TreeNode
	Current     *TreeNode
	Head        *TreeNode
	Nodes       map[int]*TreeNode
	NextIndex   int
	InputBuffer int64
	Timeout     float64
	Size        int
	Board       *Board
	Clipboard   *TreeNode
	Explorer    *Explorer
}

func (s *State) Prefs() string {
	result := "{"
	first := true
	stack := []*TreeNode{s.Root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]

		c := cur.PreferredChild
		if first {
			result = fmt.Sprintf("%s\"%d\":%d", result, cur.Index, c)
			first = false
		} else {
			result = fmt.Sprintf("%s,\"%d\":%d", result, cur.Index, c)
		}

		if len(cur.Down) == 1 {
			stack = append(stack, cur.Down[0])
		} else if len(cur.Down) > 1 {
			// go backward through array
			for i := len(cur.Down) - 1; i >= 0; i-- {
				n := cur.Down[i]
				stack = append(stack, n)
			}
		}
	}

	result += "}"
	return result
}

func (s *State) SetPreferred(index int) error {
	n := s.Nodes[index]
	cur := n
	for {
		if cur == nil {
			return fmt.Errorf("Error in indexing")
		}
		if cur.Up == nil {
			break
		}
		oldIndex := cur.Index
		cur = cur.Up
		for i, d := range cur.Down {
			if d.Index == oldIndex {
				cur.PreferredChild = i
			}
		}
	}
	return nil
}

func (s *State) ResetPrefs() {
	stack := []*TreeNode{s.Root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]

		cur.PreferredChild = 0

		if len(cur.Down) == 1 {
			stack = append(stack, cur.Down[0])
		} else if len(cur.Down) > 1 {
			// go backward through array
			for i := len(cur.Down) - 1; i >= 0; i-- {
				n := cur.Down[i]
				stack = append(stack, n)
			}
		}
	}
}

func (s *State) SetPrefs(prefs map[string]int) {
	stack := []*TreeNode{s.Root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]

		key := fmt.Sprintf("%d", cur.Index)
		p := prefs[key]

		cur.PreferredChild = p

		if len(cur.Down) == 1 {
			stack = append(stack, cur.Down[0])
		} else if len(cur.Down) > 1 {
			// go backward through array
			for i := len(cur.Down) - 1; i >= 0; i-- {
				n := cur.Down[i]
				stack = append(stack, n)
			}
		}
	}
}

func (s *State) Locate() string {
	dirs := []int{}
	c := s.Current
	for {
		myIndex := c.Index
		if c.Up == nil {
			break
		}
		u := c.Up
		for i := 0; i < len(u.Down); i++ {
			if u.Down[i].Index == myIndex {
				dirs = append(dirs, i)
			}
		}
		c = u
	}
	result := ""
	firstComma := false
	for i := len(dirs) - 1; i >= 0; i-- {
		d := dirs[i]
		if !firstComma {
			result = fmt.Sprintf("%d", d)
			firstComma = true
		} else {
			result = fmt.Sprintf("%s,%d", result, d)
		}
	}
	return result
}

func (s *State) GetNextIndex() int {
	i := s.NextIndex
	s.NextIndex++
	return i
}

func (s *State) AddFieldNode(fields map[string][]string, index int) *Diff {
	tmp := s.GetNextIndex()
	if index == -1 {
		index = tmp
	}
	n := NewTreeNode(nil, -1, index, s.Current, fields)
	s.Nodes[index] = n
	if s.Root == nil {
		s.Root = n
	} else {
		s.Current.Down = append(s.Current.Down, n)
		s.Current.PreferredChild = len(s.Current.Down) - 1
	}
	s.Current = n

	// compute diff
	diffAdd := []*StoneSet{}
	if val, ok := fields["AB"]; ok {
		add := NewCoordSet()
		for _, v := range val {
			add.Add(LettersToCoord(v))
		}
		stoneSet := NewStoneSet(add, Black)
		diffAdd = append(diffAdd, stoneSet)
	}

	if val, ok := fields["AW"]; ok {
		add := NewCoordSet()
		for _, v := range val {
			add.Add(LettersToCoord(v))
		}
		stoneSet := NewStoneSet(add, White)
		diffAdd = append(diffAdd, stoneSet)
	}

	diffRemove := []*StoneSet{}
	if val, ok := fields["AE"]; ok {
		csBlack := NewCoordSet()
		csWhite := NewCoordSet()
		for _, v := range val {
			coord := LettersToCoord(v)
			col := s.Board.Get(coord)
			if col == Black {
				csBlack.Add(coord)
			} else if col == White {
				csWhite.Add(coord)
			}
		}
		removeBlack := NewStoneSet(csBlack, Black)
		removeWhite := NewStoneSet(csWhite, White)
		diffRemove = append(diffRemove, removeBlack)
		diffRemove = append(diffRemove, removeWhite)
	}

	diff := NewDiff(diffAdd, diffRemove)
	s.Board.ApplyDiff(diff)
	s.Current.Diff = diff
	return diff
}

func (s *State) AddPassNode(col Color, fields map[string][]string, index int) {
	tmp := s.GetNextIndex()
	if index == -1 {
		index = tmp
	}
	n := NewTreeNode(nil, col, index, s.Current, fields)
	s.Nodes[index] = n
	if s.Root == nil {
		s.Root = n
	} else {
		s.Current.Down = append(s.Current.Down, n)
		s.Current.PreferredChild = len(s.Current.Down) - 1
	}
	s.Current = n
	// no need to add a diff
}

func (s *State) PushHead(x, y, col int) {
	coord := &Coord{x, y}
	if x == -1 || y == -1 {
		coord = nil
	}
	index := s.GetNextIndex()
	fields := make(map[string][]string)
	var key string
	if Color(col) == Black {
		key = "B"
	} else {
		key = "W"
	}
	value := ""
	if x != -1 {
		value = coord.ToLetters()
	}
	fields[key] = []string{value}
	n := NewTreeNode(coord, Color(col), index, s.Head, fields)
	s.Nodes[index] = n
	if len(s.Head.Down) > 0 {
		s.Head.PreferredChild++
	}
	s.Head.Down = append([]*TreeNode{n}, s.Head.Down...)

	// tracking the head or not
	tracking := false
	if s.Current == s.Head {
		tracking = true
	}

	var diff *Diff

	// if we're not tracking the head
	if !tracking {
		// save where we currently are
		save := s.Current.Index

		// goto head
		s.GotoIndex(s.Head.Index)

		// compute diff
		diff = s.Board.Move(coord, Color(col))

		// go back to saved index
		s.GotoIndex(save)
	} else {
		// do nothing if it's a pass

		// otherwise
		if x != -1 {
			// if we are tracking, just compute the diff
			diff = s.Board.Move(coord, Color(col))
		}

		// and follow along
		s.Current = n
	}

	// set new head
	s.Head = n

	// set diff
	s.Head.Diff = diff
}

func (s *State) AddNode(coord *Coord, col Color, fields map[string][]string, index int, force bool) *Diff {
	if fields == nil {
		fields = make(map[string][]string)
	}

	if !force {
		// check to see if it's already there
		for i, node := range s.Current.Down {
			coordOld := node.XY
			if coordOld != nil &&
				coord != nil &&
				coordOld.X == coord.X &&
				coordOld.Y == coord.Y &&
				node.Color == Color(col) {
				s.Current.PreferredChild = i
				s.Right()
				return s.Current.Diff
			}
		}
	}

	tmp := s.GetNextIndex()
	if index == -1 {
		index = tmp
	}
	n := NewTreeNode(coord, Color(col), index, s.Current, fields)

	s.Nodes[index] = n
	if s.Root == nil {
		s.Root = n
	} else {
		s.Current.Down = append(s.Current.Down, n)
		s.Current.PreferredChild = len(s.Current.Down) - 1
	}
	s.Current = n
	diff := s.Board.Move(coord, Color(col))
	s.Current.Diff = diff
	return diff
}

type PatternMove struct {
	Coord *Coord // nil for passes
	Color Color
}

func (s *State) AddPatternNodes(moves []*PatternMove) {
	node := s.Root
	locationSave := s.Current.Index

	for _, move := range moves {
		found := false

		for _, child := range node.Down {
			if (child.XY == nil && move.Coord == nil) ||
				(child.XY != nil && move.Coord != nil &&
					child.XY.X == move.Coord.X && child.XY.Y == move.Coord.Y &&
					child.Color == move.Color) {
				node = child
				found = true
				break
			}
		}

		if !found {
			s.GotoIndex(node.Index)

			fields := make(map[string][]string)
			key := "B"
			if move.Color == White {
				key = "W"
			}

			if move.Coord == nil {
				fields[key] = []string{""}
				s.AddPassNode(move.Color, fields, -1)
			} else {
				fields[key] = []string{move.Coord.ToLetters()}
				s.AddNode(move.Coord, move.Color, fields, -1, false)
			}
			node = s.Current
		}
	}
	s.GotoIndex(locationSave)
}

func (s *State) Cut() *Diff {
	// store the current index
	index := s.Current.Index

	// go left
	diff := s.Left()

	// find the child that matches the index to cut
	j := -1
	for i := 0; i < len(s.Current.Down); i++ {
		node := s.Current.Down[i]
		if node.Index == index {
			j = i
			break
		}
	}

	// if we didn't find anything return (shouldn't really happen)
	if j == -1 {
		return nil
	}

	// store the branch (child index j)
	branch := s.Current.Down[j]

	// cut the branch out from the children
	s.Current.Down = append(s.Current.Down[:j], s.Current.Down[j+1:]...)

	// delete all the nodes from the nodes map
	Fmap(func(n *TreeNode) {
		delete(s.Nodes, n.Index)
	}, branch)

	// adjust prefs
	if s.Current.PreferredChild >= len(s.Current.Down) {
		s.Current.PreferredChild = 0
	}

	// save the branch to the clipboard
	s.Clipboard = branch

	return diff
}

func (s *State) Left() *Diff {
	if s.Current.Up != nil {
		d := s.Current.Diff.Invert()
		s.Board.ApplyDiff(d)
		s.Current = s.Current.Up
		return d
	}
	return nil
}

func (s *State) Right() *Diff {
	if len(s.Current.Down) > 0 {
		index := s.Current.PreferredChild
		s.Current = s.Current.Down[index]
		d := s.Current.Diff
		s.Board.ApplyDiff(d)
		return d
	}
	return nil
}

func (s *State) GotoIndex(index int) error {
	err := s.SetPreferred(index)
	if err != nil {
		return err
	}
	s.Rewind()
	for {
		if s.Current.Index == index {
			break
		}
		s.Right()
	}
	//s.Current = s.Nodes[index]
	return nil
}

func (s *State) Rewind() {
	s.Current = s.Root
	s.Board.Clear()
	s.Board.ApplyDiff(s.Current.Diff)
}

func (s *State) FastForward() {
	for {
		if len(s.Current.Down) == 0 {
			break
		}
		s.Right()
		/*
		   index := s.Current.PreferredChild
		   s.Current = s.Current.Down[index]
		*/
	}
}

func (s *State) Up() {
	if len(s.Current.Down) == 0 {
		return
	}
	c := s.Current.PreferredChild
	mod := len(s.Current.Down)
	s.Current.PreferredChild = (((c - 1) % mod) + mod) % mod
}

func (s *State) Down() {
	if len(s.Current.Down) == 0 {
		return
	}
	c := s.Current.PreferredChild
	mod := len(s.Current.Down)
	s.Current.PreferredChild = (((c + 1) % mod) + mod) % mod
}

func (s *State) GotoCoord(x, y int) {
	cur := s.Current
	// look forward
	for {
		if cur.XY != nil && cur.XY.X == x && cur.XY.Y == y {
			s.GotoIndex(cur.Index)
			return
		}
		if len(cur.Down) == 0 {
			break
		}
		cur = cur.Down[cur.PreferredChild]
	}

	cur = s.Current
	// look backward
	for {
		if cur.XY != nil && cur.XY.X == x && cur.XY.Y == y {
			s.GotoIndex(cur.Index)
			return
		}
		if cur.Up == nil {
			break
		}
		cur = cur.Up
	}

	// didn't find anything
	// do nothing
}

func (s *State) GenerateMarks() *Marks {
	marks := &Marks{}
	if s.Current.XY != nil {
		marks.Current = s.Current.XY
	}
	if trs, ok := s.Current.Fields["TR"]; ok {
		cs := NewCoordSet()
		for _, tr := range trs {
			c := LettersToCoord(tr)
			cs.Add(c)
		}
		marks.Triangles = cs.List()
	}
	if sqs, ok := s.Current.Fields["SQ"]; ok {
		cs := NewCoordSet()
		for _, sq := range sqs {
			c := LettersToCoord(sq)
			cs.Add(c)
		}
		marks.Squares = cs.List()
	}
	if lbs, ok := s.Current.Fields["LB"]; ok {
		labels := []*Label{}
		for _, lb := range lbs {
			spl := strings.Split(lb, ":")
			c := LettersToCoord(spl[0])
			text := spl[1]
			label := &Label{c, text}
			labels = append(labels, label)
		}
		marks.Labels = labels
	}

	if pxs, ok := s.Current.Fields["PX"]; ok {
		pens := []*Pen{}
		for _, px := range pxs {
			spl := strings.Split(px, ":")
			if len(spl) != 5 {
				continue
			}
			x0, err := strconv.ParseFloat(spl[0], 64)
			y0, err := strconv.ParseFloat(spl[1], 64)
			x1, err := strconv.ParseFloat(spl[2], 64)
			y1, err := strconv.ParseFloat(spl[3], 64)
			if err != nil {
				continue
			}
			pen := &Pen{x0, y0, x1, y1, spl[4]}
			pens = append(pens, pen)
		}
		marks.Pens = pens
	}
	return marks
}

func (s *State) GenerateMetadata() *Metadata {
	m := &Metadata{
		Size:   s.Size,
		Fields: s.Root.Fields,
	}
	return m
}

func (s *State) GenerateComments() []string {
	cmts := []string{}
	if c, ok := s.Current.Fields["C"]; ok {
		cmts = c
	}
	return cmts
}

func (s *State) GenerateFullFrame(t TreeJSONType) *Frame {
	frame := s.Board.CurrentFrame()
	frame.Marks = s.GenerateMarks()
	frame.Metadata = s.GenerateMetadata()
	frame.Comments = s.GenerateComments()
	frame.TreeJSON = s.CreateTreeJSON(t)
	return frame
}

// full frame or diff frame
// explorer:
//	- change current
//	- change preferred path
//	- change tree

/*
func (s *State) GenerateFrame(t ExplorerType, f FrameType, d *Diff) *Frame {
	marks := s.GenerateMarks()
	comments := s.GenerateMetadata()
	metadata := s.GenerateMetadata()

	switch t {
	case CurrentOnly:
		explorer := &Explorer{}
		explorer.Current = s.Explorer.Current.Coord

	case CurrentAndPreferred:


	case Full:
	}
}
*/

// see addevent.go
func (s *State) AddEvent(evt *EventJSON) (*Frame, error) {
	switch evt.Event {
	case "add_stone":
		return s.HandleAddStone(evt)
	case "pass":
		return s.HandlePass(evt)
	case "remove_stone":
		return s.HandleRemoveStone(evt)
	case "triangle":
		return s.HandleAddTriangle(evt)
	case "square":
		return s.HandleAddSquare(evt)
	case "letter":
		return s.HandleAddLetter(evt)
	case "number":
		return s.HandleAddNumber(evt)
	case "remove_mark":
		return s.HandleRemoveMark(evt)
	case "cut":
		return s.HandleCut(evt)
	case "left":
		return s.HandleLeft()
	case "right":
		return s.HandleRight()
	case "up":
		return s.HandleUp()
	case "down":
		return s.HandleDown()
	case "rewind":
		return s.HandleRewind()
	case "fastforward":
		return s.HandleFastForward()
	case "goto_grid":
		return s.HandleGotoGrid(evt)
	case "goto_coord":
		return s.HandleGotoCoord(evt)
	case "comment":
		return s.HandleComment(evt)
	case "draw":
		return s.HandleDraw(evt)
	case "erase_pen":
		return s.HandleErasePen()
	case "copy":
		return s.HandleCopy()
	case "clipboard":
		return s.HandleClipboard()
	}
	return nil, nil
}

func (s *State) ToSGF(indexes bool) string {
	result := "("
	stack := []interface{}{s.Root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]
		if str, ok := cur.(string); ok {
			result += str
			continue
		}
		node := cur.(*TreeNode)
		result += ";"
		// throw in other fields
		for key, multifield := range node.Fields {
			if key == "IX" {
				continue
			}
			result += fmt.Sprintf("%s", key)
			for _, fieldValue := range multifield {
				m := strings.ReplaceAll(fieldValue, "]", "\\]")
				result += fmt.Sprintf("[%s]", m)
			}

		}

		if indexes {
			result += fmt.Sprintf("IX[%d]", node.Index)
		}

		if len(node.Down) == 1 {
			stack = append(stack, node.Down[0])
		} else if len(node.Down) > 1 {
			// go backward through array
			for i := len(node.Down) - 1; i >= 0; i-- {
				n := node.Down[i]
				stack = append(stack, ")")
				stack = append(stack, n)
				stack = append(stack, "(")
			}
		}
	}

	result += ")"
	return result
}

func FromSGF(data string) (*State, error) {
	p := NewParser(data)
	root, err := p.Parse()
	if err != nil {
		return nil, err
	}

	var size int64 = 19
	if _, ok := root.Fields["SZ"]; ok {
		size_field := root.Fields["SZ"]
		if len(size_field) != 1 {
			return nil, fmt.Errorf("SZ cannot be a multifield")
		}
		size, err = strconv.ParseInt(size_field[0], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	state := NewState(int(size), false)
	stack := []interface{}{root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]
		if _, ok := cur.(string); ok {
			state.Left()
		} else {
			node := cur.(*SGFNode)

			index := -1
			if indexes, ok := node.Fields["IX"]; ok {
				if len(indexes) > 0 {
					_index, err := strconv.ParseInt(indexes[0], 10, 64)
					index = int(_index)
					if err != nil {
						index = -1
					}
				}
			}

			// refuse to process sgfs with a suicide move
			if node.Coord() != nil && !state.Board.Legal(node.Coord(), node.Color()) {
				return nil, fmt.Errorf("suicide moves are not currently supported")
			}

			if node.IsPass() {
				state.AddPassNode(node.Color(), node.Fields, index)
			} else if node.IsMove() {
				state.AddNode(node.Coord(), node.Color(), node.Fields, index, true)
			} else {
				state.AddFieldNode(node.Fields, index)
			}
			for i := len(node.Down) - 1; i >= 0; i-- {
				stack = append(stack, "<")
				stack = append(stack, node.Down[i])
			}
			state.Head = state.Current
		}
	}
	state.Rewind()
	state.ResetPrefs()
	return state, nil
}

func (s *State) InitData(event string) *EventJSON {
	sgf := s.ToSGF(true)
	encoded := base64.StdEncoding.EncodeToString([]byte(sgf))
	loc := s.Locate()
	prefs := s.Prefs()
	value := fmt.Sprintf("{\"sgf\":\"%s\", \"loc\":\"%s\", \"prefs\":%s, \"buffer\":%d, \"next_index\":%d}", encoded, loc, prefs, s.InputBuffer, s.NextIndex)
	evt := &EventJSON{event, value, 0, ""}
	return evt

	//return []byte(fmt.Sprintf("{\"event\":\"%s\",\"value\":%s}", event, value))

}

func NewState(size int, initRoot bool) *State {
	nodes := make(map[int]*TreeNode)
	var root *TreeNode
	root = nil
	index := 0
	if initRoot {
		fields := map[string][]string{}
		fields["GM"] = []string{"1"}
		fields["FF"] = []string{"4"}
		fields["CA"] = []string{"UTF-8"}
		fields["SZ"] = []string{fmt.Sprintf("%d", size)}
		fields["PB"] = []string{"Black"}
		fields["PW"] = []string{"White"}
		fields["RU"] = []string{"Japanese"}
		fields["KM"] = []string{"6.5"}

		// coord, color, index, up, fields
		root = NewTreeNode(nil, 0, 0, nil, fields)
		nodes[0] = root
		index = 1
	}
	board := NewBoard(size)
	// default input buffer of 250
	// default room timeout of 86400
	return &State{root, root, root, nodes, index, 250, 86400, size, board, nil, nil}
}
