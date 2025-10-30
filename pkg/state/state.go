/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/jarednogo/board/pkg/core"
)

const Letters = "ABCDEFGHIJKLNMOPQRSTUVWXYZ"

// as a rule, anything that would need to get sent to new connections
// should be stored here
type State struct {
	root        *core.TreeNode
	current     *core.TreeNode
	head        *core.TreeNode
	nodes       map[int]*core.TreeNode
	nextIndex   int
	inputBuffer int64
	timeout     float64
	size        int
	board       *core.Board
	clipboard   *core.TreeNode
	markedDead  core.CoordSet
	markedDame  core.CoordSet
}

func (s *State) HeadColor() core.Color {
	return s.head.Color
}

func (s *State) Prefs() map[string]int {
	prefs := make(map[string]int)
	for index, node := range s.nodes {
		key := fmt.Sprintf("%d", index)
		prefs[key] = node.PreferredChild
	}
	return prefs
}

func (s *State) SetPreferred(index int) error {
	n := s.nodes[index]
	cur := n
	for {
		if cur == nil {
			return fmt.Errorf("error in indexing")
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

// ResetPrefs sets all prefs to 0
func (s *State) ResetPrefs() {
	for _, n := range s.nodes {
		n.PreferredChild = 0
	}
}

// SetPrefs takes a map and sets
func (s *State) SetPrefs(prefs map[string]int) {
	for _, n := range s.nodes {
		key := fmt.Sprintf("%d", n.Index)
		p := prefs[key]
		n.PreferredChild = p
	}
}

func (s *State) Locate() string {
	dirs := []int{}
	c := s.current
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

func (s *State) SetNextIndex(i int) {
	s.nextIndex = i
}

func (s *State) GetNextIndex() int {
	i := s.nextIndex
	s.nextIndex++
	return i
}

func (s *State) GetInputBuffer() int64 {
	return s.inputBuffer
}

func (s *State) SetInputBuffer(i int64) {
	s.inputBuffer = i
}

func (s *State) GetTimeout() float64 {
	return s.timeout
}

func (s *State) SetTimeout(f float64) {
	s.timeout = f
}

func (s *State) Size() int {
	return s.size
}

func (s *State) Board() *core.Board {
	return s.board
}

func (s *State) Current() *core.TreeNode {
	return s.current
}

func (s *State) Root() *core.TreeNode {
	return s.root
}

func (s *State) Head() *core.TreeNode {
	return s.head
}

func (s *State) Nodes() map[int]*core.TreeNode {
	return s.nodes
}

func (s *State) AnyMove() {
	s.markedDead = core.NewCoordSet()
	s.markedDame = core.NewCoordSet()
}

func (s *State) AddFieldNode(fields map[string][]string, index int) *core.Diff {
	s.AnyMove()
	tmp := s.GetNextIndex()
	if index == -1 {
		index = tmp
	}
	n := core.NewTreeNode(nil, core.NoColor, index, s.current, fields)
	s.nodes[index] = n
	if s.root == nil {
		s.root = n
	} else {
		s.current.Down = append(s.current.Down, n)
		s.current.PreferredChild = len(s.current.Down) - 1
	}
	s.current = n

	// compute diff
	diff := s.ComputeDiffSetup(index)
	s.board.ApplyDiff(diff)
	s.current.SetDiff(diff)
	return diff
}

func (s *State) AddPassNode(col core.Color, fields map[string][]string, index int) {
	s.AnyMove()
	tmp := s.GetNextIndex()
	if index == -1 {
		index = tmp
	}
	n := core.NewTreeNode(nil, col, index, s.current, fields)
	s.nodes[index] = n
	if s.root == nil {
		s.root = n
	} else {
		s.current.Down = append(s.current.Down, n)
		s.current.PreferredChild = len(s.current.Down) - 1
	}
	s.current = n
	// no need to add a diff
	// but actually, SetDiff also sets the score
	s.current.SetDiff(nil)
}

func (s *State) PushHead(x, y int, col core.Color) {
	coord := &core.Coord{X: x, Y: y}
	if x == -1 || y == -1 {
		coord = nil
	}
	index := s.GetNextIndex()
	fields := make(map[string][]string)
	var key string
	if col == core.Black {
		key = "B"
	} else {
		key = "W"
	}
	value := ""
	if x != -1 {
		value = coord.ToLetters()
	}
	fields[key] = []string{value}
	n := core.NewTreeNode(coord, col, index, s.head, fields)
	s.nodes[index] = n
	if len(s.head.Down) > 0 {
		s.head.PreferredChild++
	}
	s.head.Down = append([]*core.TreeNode{n}, s.head.Down...)

	// tracking the head or not
	tracking := s.current == s.head

	var diff *core.Diff

	// if we're not tracking the head
	if !tracking {
		// save where we currently are
		save := s.current.Index

		// goto head
		s.GotoIndex(s.head.Index) //nolint: errcheck

		// compute diff
		diff = s.board.Move(coord, col)

		// go back to saved index
		s.GotoIndex(save) //nolint: errcheck
	} else {
		// do nothing if it's a pass

		// otherwise
		if x != -1 {
			// if we are tracking, just compute the diff
			diff = s.board.Move(coord, col)
		}

		// and follow along
		s.current = n
	}

	// set new head
	s.head = n

	// set diff
	s.head.SetDiff(diff)
}

func (s *State) AddNode(coord *core.Coord, col core.Color, fields map[string][]string, index int, force bool) *core.Diff {
	s.AnyMove()
	if fields == nil {
		fields = make(map[string][]string)
	}

	if !force {
		// check to see if it's already there
		for i, node := range s.current.Down {
			coordOld := node.XY
			if coordOld != nil &&
				coord != nil &&
				coordOld.X == coord.X &&
				coordOld.Y == coord.Y &&
				node.Color == core.Color(col) {
				s.current.PreferredChild = i
				s.Right()
				return s.current.Diff
			}
		}
	}

	tmp := s.GetNextIndex()
	if index == -1 {
		index = tmp
	}
	n := core.NewTreeNode(coord, core.Color(col), index, s.current, fields)

	s.nodes[index] = n
	if s.root == nil {
		s.root = n
	} else {
		s.current.Down = append(s.current.Down, n)
		s.current.PreferredChild = len(s.current.Down) - 1
	}
	s.current = n
	diff := s.board.Move(coord, core.Color(col))
	s.current.SetDiff(diff)
	return diff
}

func (s *State) AddPatternNodes(moves []*core.PatternMove) {
	node := s.root
	locationSave := s.current.Index

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
			s.GotoIndex(node.Index) //nolint: errcheck

			fields := make(map[string][]string)
			key := "B"
			if move.Color == core.White {
				key = "W"
			}

			if move.Coord == nil {
				fields[key] = []string{""}
				s.AddPassNode(move.Color, fields, -1)
			} else {
				fields[key] = []string{move.Coord.ToLetters()}
				s.AddNode(move.Coord, move.Color, fields, -1, false)
			}
			node = s.current
		}
	}
	s.GotoIndex(locationSave) //nolint: errcheck
}

// SmartGraft doesn't duplicate existing moves
func (s *State) SmartGraft(parentIndex int, moves []*core.PatternMove) {
	parent := s.nodes[parentIndex]
	savedPrefs := make(map[int]int)
	save := s.current.Index

	var graft *core.TreeNode
	up := parent

	for _, move := range moves {

		// go to the parent
		s.GotoIndex(up.Index) //nolint: errcheck

		// save the preferences on each node that already exists
		savedPrefs[up.Index] = up.PreferredChild

		// if the move exists in a child node, then follow it
		if i, ok := s.current.HasChild(move.Coord, move.Color); ok {
			up = s.nodes[i]
			continue
		}

		// once we get here we are adding new nodes

		// each node needs an index
		index := s.GetNextIndex()

		// each node needs either B[] or W[] field
		fields := make(map[string][]string)
		var key string
		if move.Color == core.Black {
			key = "B"
		} else {
			key = "W"
		}
		fields[key] = []string{move.Coord.ToLetters()}

		// create the node, up is the parent of the new node
		node := core.NewTreeNode(move.Coord, move.Color, index, up, fields)

		// keep track of the first new node
		if graft == nil {
			graft = node
		}

		// add the node to the state's node map
		s.nodes[index] = node

		// follow along so we can set child nodes
		up.Down = append(up.Down, node)

		// calculate the diff
		diff := s.board.Move(move.Coord, move.Color)
		node.SetDiff(diff)

		// set the new parent for the next node
		up = node
	}

	// cleanup

	// (this is only necessary if we added something)
	if graft != nil {
		graft.RecomputeDepth()
	}

	s.GotoIndex(save) //nolint: errcheck
	for index, pref := range savedPrefs {
		s.nodes[index].PreferredChild = pref
	}

}

// Graft may duplicate existing moves
func (s *State) Graft(parentIndex int, moves []*core.PatternMove) {
	parent := s.nodes[parentIndex]
	savedPref := parent.PreferredChild
	save := s.current.Index

	var graft *core.TreeNode
	up := parent

	for _, move := range moves {

		// go to the parent
		s.GotoIndex(up.Index) //nolint: errcheck

		// each node needs an index
		index := s.GetNextIndex()

		// each node needs either B[] or W[] field
		fields := make(map[string][]string)
		var key string
		if move.Color == core.Black {
			key = "B"
		} else {
			key = "W"
		}
		fields[key] = []string{move.Coord.ToLetters()}

		// create the node, up is the parent of the new node
		node := core.NewTreeNode(move.Coord, move.Color, index, up, fields)

		// keep track of the first node
		if graft == nil {
			graft = node
		}

		// add the node to the state's node map
		s.nodes[index] = node

		// follow along so we can set child nodes
		up.Down = append(up.Down, node)

		// calculate the diff
		diff := s.board.Move(move.Coord, move.Color)
		node.SetDiff(diff)

		// set the new parent for the next node
		up = node
	}

	// cleanup
	graft.RecomputeDepth()
	s.GotoIndex(save) //nolint: errcheck
	parent.PreferredChild = savedPref
}

func (s *State) Cut() *core.Diff {
	s.AnyMove()
	// store the current index
	index := s.current.Index

	// go left
	diff := s.Left()

	// find the child that matches the index to cut
	j := -1
	for i := 0; i < len(s.current.Down); i++ {
		node := s.current.Down[i]
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
	branch := s.current.Down[j]

	// cut the branch out from the children
	s.current.Down = append(s.current.Down[:j], s.current.Down[j+1:]...)

	// delete all the nodes from the nodes map
	core.Fmap(func(n *core.TreeNode) {
		delete(s.nodes, n.Index)
	}, branch)

	// adjust prefs
	if s.current.PreferredChild >= len(s.current.Down) {
		s.current.PreferredChild = 0
	}

	// save the branch to the clipboard
	s.clipboard = branch

	return diff
}

func (s *State) Left() *core.Diff {
	s.AnyMove()
	if s.current.Up != nil {
		d := s.current.Diff.Invert()
		s.board.ApplyDiff(d)
		s.current = s.current.Up
		return d
	}
	return nil
}

func (s *State) Right() *core.Diff {
	s.AnyMove()
	if len(s.current.Down) > 0 {
		index := s.current.PreferredChild
		s.current = s.current.Down[index]
		d := s.current.Diff
		s.board.ApplyDiff(d)
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
	last := s.current.Index
	for s.current.Index != index {
		s.Right()
		// to prevent infinite loops
		if s.current.Index == last {
			break
		}
	}
	//s.Current = s.Nodes[index]
	return nil
}

func (s *State) Rewind() {
	s.AnyMove()
	s.current = s.root
	s.board.Clear()
	s.board.ApplyDiff(s.current.Diff)
}

func (s *State) FastForward() {
	for len(s.current.Down) != 0 {
		s.Right()
	}
}

func (s *State) Up() {
	if len(s.current.Down) == 0 {
		return
	}
	c := s.current.PreferredChild
	mod := len(s.current.Down)
	s.current.PreferredChild = (((c - 1) % mod) + mod) % mod
}

func (s *State) Down() {
	if len(s.current.Down) == 0 {
		return
	}
	c := s.current.PreferredChild
	mod := len(s.current.Down)
	s.current.PreferredChild = (((c + 1) % mod) + mod) % mod
}

func (s *State) GotoCoord(x, y int) {
	cur := s.current
	// look forward
	for {
		if cur.XY != nil && cur.XY.X == x && cur.XY.Y == y {
			s.GotoIndex(cur.Index) //nolint: errcheck
			return
		}
		if len(cur.Down) == 0 {
			break
		}
		cur = cur.Down[cur.PreferredChild]
	}

	cur = s.current
	// look backward
	for {
		if cur.XY != nil && cur.XY.X == x && cur.XY.Y == y {
			s.GotoIndex(cur.Index) //nolint: errcheck
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

func (s *State) GenerateMarks() *core.Marks {
	marks := &core.Marks{}
	if s.current.XY != nil {
		marks.Current = s.current.XY
	}
	if trs, ok := s.current.Fields["TR"]; ok {
		cs := core.NewCoordSet()
		for _, tr := range trs {
			c := core.LettersToCoord(tr)
			cs.Add(c)
		}
		marks.Triangles = cs.List()
	}
	if sqs, ok := s.current.Fields["SQ"]; ok {
		cs := core.NewCoordSet()
		for _, sq := range sqs {
			c := core.LettersToCoord(sq)
			cs.Add(c)
		}
		marks.Squares = cs.List()
	}
	if lbs, ok := s.current.Fields["LB"]; ok {
		labels := []*core.Label{}
		for _, lb := range lbs {
			spl := strings.Split(lb, ":")
			c := core.LettersToCoord(spl[0])
			text := spl[1]
			label := &core.Label{Coord: c, Text: text}
			labels = append(labels, label)
		}
		marks.Labels = labels
	}

	if pxs, ok := s.current.Fields["PX"]; ok {
		pens := []*core.Pen{}
		for _, px := range pxs {
			spl := strings.Split(px, ":")
			if len(spl) != 5 {
				continue
			}
			x0, err0 := strconv.ParseFloat(spl[0], 64)
			y0, err1 := strconv.ParseFloat(spl[1], 64)
			x1, err2 := strconv.ParseFloat(spl[2], 64)
			y1, err3 := strconv.ParseFloat(spl[3], 64)
			hasErr := err0 != nil || err1 != nil || err2 != nil || err3 != nil
			if hasErr {
				continue
			}
			pen := &core.Pen{X0: x0, Y0: y0, X1: x1, Y1: y1, Color: spl[4]}
			pens = append(pens, pen)
		}
		marks.Pens = pens
	}
	return marks
}

func (s *State) GenerateMetadata() *core.Metadata {
	m := &core.Metadata{
		Size:   s.size,
		Fields: s.root.Fields,
	}
	return m
}

func (s *State) GenerateComments() []string {
	cmts := []string{}
	if c, ok := s.current.Fields["C"]; ok {
		cmts = c
	}
	return cmts
}

func (s *State) GenerateFullFrame(t core.TreeJSONType) *core.Frame {
	frame := &core.Frame{}
	frame.Type = core.FullFrame
	frame.Diff = s.board.CurrentDiff()
	frame.Marks = s.GenerateMarks()
	frame.Metadata = s.GenerateMetadata()
	frame.Comments = s.GenerateComments()
	frame.TreeJSON = s.CreateTreeJSON(t)
	frame.BlackCaps = s.current.BlackCaps
	frame.WhiteCaps = s.current.WhiteCaps
	return frame
}

func (s *State) ComputeDiffMove(i int) *core.Diff {
	save := s.current.Index
	n, ok := s.nodes[i]
	if !ok {
		return nil
	}
	// only for the purposes of checking the board (for removal of stones)
	if n.Up != nil {
		s.GotoIndex(n.Up.Index) //nolint: errcheck
	} else {
		s.board.Clear()
	}

	// at the end, go back to the saved index
	defer s.GotoIndex(save) //nolint: errcheck
	if n.XY == nil {
		return nil
	}

	diff := s.board.Move(n.XY, n.Color)
	return diff
}

func (s *State) ComputeDiffSetup(i int) *core.Diff {
	save := s.current.Index
	n, ok := s.nodes[i]
	if !ok {
		return nil
	}

	// only for the purposes of checking the board (for removal of stones)
	if n.Up != nil {
		s.GotoIndex(n.Up.Index) //nolint: errcheck
	} else {
		s.board.Clear()
	}

	// at the end, go back to the saved index
	defer s.GotoIndex(save) //nolint: errcheck

	// find the black stones to add
	diffAdd := []*core.StoneSet{}
	if val, ok := n.Fields["AB"]; ok {
		add := core.NewCoordSet()
		for _, v := range val {
			add.Add(core.LettersToCoord(v))
		}
		stoneSet := core.NewStoneSet(add, core.Black)
		diffAdd = append(diffAdd, stoneSet)
	}

	// find the white stones to add
	if val, ok := n.Fields["AW"]; ok {
		add := core.NewCoordSet()
		for _, v := range val {
			add.Add(core.LettersToCoord(v))
		}
		stoneSet := core.NewStoneSet(add, core.White)
		diffAdd = append(diffAdd, stoneSet)
	}

	// find the stones to remove
	diffRemove := []*core.StoneSet{}
	if val, ok := n.Fields["AE"]; ok {
		csBlack := core.NewCoordSet()
		csWhite := core.NewCoordSet()
		for _, v := range val {
			coord := core.LettersToCoord(v)
			col := s.board.Get(coord)
			switch col {
			case core.Black:
				csBlack.Add(coord)
			case core.White:
				csWhite.Add(coord)
			}
		}
		removeBlack := core.NewStoneSet(csBlack, core.Black)
		removeWhite := core.NewStoneSet(csWhite, core.White)
		diffRemove = append(diffRemove, removeBlack)
		diffRemove = append(diffRemove, removeWhite)
	}
	diff := core.NewDiff(diffAdd, diffRemove)
	return diff
}

// see addevent.go
func (s *State) AddEvent(evt *core.EventJSON) (*core.Frame, error) {
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
		return s.HandleCut()
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
	case "graft":
		return s.HandleGraft(evt)
	case "score":
		return s.HandleScore()
	case "markdead":
		return s.HandleMarkDead(evt)
	}
	return nil, nil
}

func (s *State) ToSGF(indexes bool) string {
	result := "("
	stack := []interface{}{s.root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]
		if str, ok := cur.(string); ok {
			result += str
			continue
		}
		node := cur.(*core.TreeNode)
		result += ";"
		// throw in other fields
		for key, multifield := range node.Fields {
			if key == "IX" {
				continue
			}
			result += key
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
	p := core.NewParser(data)
	root, err := p.Parse()
	if err != nil {
		return nil, err
	}

	var size int64 = 19
	if _, ok := root.Fields["SZ"]; ok {
		sizeField := root.Fields["SZ"]
		if len(sizeField) != 1 {
			return nil, fmt.Errorf("SZ cannot be a multifield")
		}
		size, err = strconv.ParseInt(sizeField[0], 10, 64)
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
			node := cur.(*core.SGFNode)

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
			if node.Coord() != nil && !state.board.Legal(node.Coord(), node.Color()) {
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
			// TODO: this might be wrong in some cases
			state.head = state.current
		}
	}
	state.Rewind()
	state.ResetPrefs()
	return state, nil
}

type StateJSON struct {
	SGF       string         `json:"sgf"`
	Location  string         `json:"loc"`
	Prefs     map[string]int `json:"prefs"`
	Buffer    int64          `json:"buffer"`
	NextIndex int            `json:"next_index"`
}

func (s *State) CreateStateJSON() *StateJSON {
	sgf := s.ToSGF(true)
	encoded := base64.StdEncoding.EncodeToString([]byte(sgf))
	loc := s.Locate()
	prefs := s.Prefs()
	stateStruct := &StateJSON{
		SGF:       encoded,
		Location:  loc,
		Prefs:     prefs,
		Buffer:    s.inputBuffer,
		NextIndex: s.nextIndex,
	}
	return stateStruct
	//value := fmt.Sprintf("{\"sgf\":\"%s\", \"loc\":\"%s\", \"prefs\":%s, \"buffer\":%d, \"next_index\":%d}", encoded, loc, prefs, s.InputBuffer, s.NextIndex)
	//evt := &core.EventJSON{"init", value, 0, ""}
	//return evt

	//return []byte(fmt.Sprintf("{\"event\":\"%s\",\"value\":%s}", event, value))

}

func NewState(size int, initRoot bool) *State {
	nodes := make(map[int]*core.TreeNode)
	var root *core.TreeNode
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
		root = core.NewTreeNode(nil, core.NoColor, 0, nil, fields)
		nodes[0] = root
		index = 1
	}
	board := core.NewBoard(size)
	// default input buffer of 250
	// default timeout of 86400
	return &State{root, root, root, nodes, index, 250, 86400, size, board, nil, core.NewCoordSet(), core.NewCoordSet()}
}

func (s *State) CreateTreeJSON(t core.TreeJSONType) *core.TreeJSON {
	// only really used when we have a partial tree
	up := 0
	root := 0

	// nodes
	var nodes map[int]*core.NodeJSON

	if t >= core.PartialNodes {
		nodes = make(map[int]*core.NodeJSON)
		var start *core.TreeNode
		// we can choose to send the full or just a partial tree
		// based on which node we start on

		switch t {
		case core.PartialNodes:
			start = s.current
			up = start.Up.Index
			root = start.Index
		case core.Full:
			start = s.root
		}
		core.Fmap(func(n *core.TreeNode) {
			down := []int{}
			for _, c := range n.Down {
				down = append(down, c.Index)
			}
			nodes[n.Index] = &core.NodeJSON{
				Color: n.Color,
				Down:  down,
				Depth: n.Depth,
			}

		}, start)
	}

	// preferred
	var preferred []int
	if t >= core.CurrentAndPreferred {
		node := s.root
		preferred = []int{node.Index}
		for len(node.Down) != 0 {
			node = node.Down[node.PreferredChild]
			preferred = append(preferred, node.Index)
		}
	}

	return &core.TreeJSON{
		Nodes:     nodes,
		Current:   s.current.Index,
		Preferred: preferred,
		Depth:     s.root.MaxDepth(),
		Up:        up,
		Root:      root,
	}
}
