/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core

/*
reference: https://www.red-bean.com/sgf/user_guide/index.html#move_vs_place
It is good style (and is required since FF[4]) to distinguish between a move and the position arrived at by this move.

Therefore it's illegal to mix setup properties and move properties within the same node.

full list of properties: https://www.red-bean.com/sgf/proplist_t.html
B and W are move properties
AB, AE, and AW are setup properties
*/

import (
	"fmt"

	"github.com/jarednogo/board/pkg/core/color"
)

type Group struct {
	Coords CoordSet
	Libs   CoordSet
	Color  color.Color
}

func (g *Group) String() string {
	return fmt.Sprintf("(%v, %v)", g.Coords, g.Color)
}

func NewGroup(coords CoordSet, libs CoordSet, col color.Color) *Group {
	if coords == nil {
		coords = NewCoordSet()
	}
	if libs == nil {
		libs = NewCoordSet()
	}
	return &Group{
		Coords: coords,
		Libs:   libs,
		Color:  col,
	}
}

type Board struct {
	Size   int
	Points [][]color.Color
}

func NewBoard(size int) *Board {
	points := [][]color.Color{}
	for i := 0; i < size; i++ {
		row := make([]color.Color, size)
		points = append(points, row)
	}
	return &Board{
		Size:   size,
		Points: points,
	}
}

func (b *Board) String() string {
	result := ""
	for _, row := range b.Points {
		for _, c := range row {
			result += fmt.Sprintf("%v ", c)
		}
		result += "\n"
	}
	return result
}

func (b *Board) Clear() {
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			b.Points[i][j] = color.Empty
		}
	}
}

func (b *Board) Copy() *Board {
	c := NewBoard(b.Size)
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			c.Points[i][j] = b.Points[i][j]
		}
	}
	return c
}

func (b *Board) Set(c *Coord, col color.Color) {
	b.Points[c.Y][c.X] = col
}

func (b *Board) Get(c *Coord) color.Color {
	if c == nil {
		return color.Empty
	}
	if c.Y >= b.Size || c.X >= b.Size || c.Y < 0 || c.X < 0 {
		return color.Empty
	}
	return b.Points[c.Y][c.X]
}

func (b *Board) SetMany(cs []*Coord, col color.Color) {
	for _, c := range cs {
		b.Set(c, col)
	}
}

func (b *Board) Neighbors(c *Coord) CoordSet {
	nbs := NewCoordSet()
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if (x != 0 && y != 0) || (x == 0 && y == 0) {
				continue
			}
			newX := c.X + x
			newY := c.Y + y
			if newX < 0 || newY < 0 {
				continue
			}
			if newX >= b.Size || newY >= b.Size {
				continue
			}
			nbs.Add(NewCoord(newX, newY))
		}
	}
	return nbs
}

func (b *Board) FindGroup(start *Coord) *Group {
	// get the color of the starting point
	col := b.Get(start)

	// if it's empty, return empty group
	if col == color.Empty {
		return NewGroup(nil, nil, color.Empty)
	}

	// initiate the stack
	stack := []*Coord{start}

	// keep track of liberties as we go
	// map so we don't double count
	libs := NewCoordSet()

	// initiate elements
	elts := NewCoordSet()

	// start DFS
	for len(stack) > 0 {
		// pop off the stack
		point := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// add to elements
		elts.Add(point)

		// compute neighbors
		nbs := b.Neighbors(point)
		for _, nb := range nbs {
			// if it's the right color
			// and we haven't visited it yet
			// add to the stack
			if col.Equal(b.Get(nb)) && !elts.Has(nb) {
				stack = append(stack, nb)
			} else if b.Get(nb) == color.Empty {
				libs.Add(nb)
			}
		}
	}
	return NewGroup(elts, libs, col)
}

func (b *Board) Groups() []*Group {
	// keep track of which points we've covered
	check := make(map[[2]int]bool)

	groups := []*Group{}

	// go through the whole board
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			coord := NewCoord(i, j)
			// if we haven't checked it yet and there's a stone here
			if !check[[2]int{i, j}] && (b.Get(coord) == color.Black || b.Get(coord) == color.White) {
				// find the group it's part of
				gp := b.FindGroup(coord)
				for _, c := range gp.Coords {
					// check off everything in the group
					check[[2]int{c.X, c.Y}] = true
				}
				// add to the list of groups
				groups = append(groups, gp)
			}
		}
	}
	return groups
}

func (b *Board) Legal(start *Coord, col color.Color) bool {
	// if there's already a stone there, it's illegal
	if b.Get(start) != color.Empty {
		return false
		// not legal
	}

	// this should be undone at the end
	b.Set(start, col)
	defer b.Set(start, color.Empty)

	// if it has >0 libs, it's legal
	gp := b.FindGroup(start)
	if len(gp.Libs) > 0 {
		return true
	}

	// check for any groups of opposite color with 0 libs
	// only check neighboring area for optimization
	nbs := b.Neighbors(start)
	for _, nb := range nbs {
		if b.Get(nb) == color.Empty {
			// unreachable
			continue
		}
		gp := b.FindGroup(nb)
		if len(gp.Libs) == 0 && gp.Color == col.Opposite() {
			// if we killed something, it's legal
			return true
		}
	}

	// if we have 0 libs and we didn't kill anything
	// it's a suicide move (and not legal)
	return false
}

func (b *Board) WouldKill(start *Coord, col color.Color) *StoneSet {
	// we pretend a stone of color Opposite(col) was just played at start
	a := b.Get(start)
	b.Set(start, col.Opposite())
	defer b.Set(start, a)
	dead := NewCoordSet()
	for _, nb := range b.Neighbors(start) {
		// if we've already marked the stone dead
		// or it's the wrong color
		// just move on
		if dead.Has(nb) || b.Get(nb) != col {
			continue
		}
		// find the group
		gp := b.FindGroup(nb)
		// if it's dead, add each to the list
		if len(gp.Libs) == 0 {
			for _, coord := range gp.Coords {
				dead.Add(coord)
			}
		}
	}
	return NewStoneSet(dead, col)
}

func (b *Board) RemoveDead(start *Coord, col color.Color) *StoneSet {
	w := b.WouldKill(start, col)
	b.SetMany(w.Coords, color.Empty)
	return w
}

func (b *Board) Move(start *Coord, col color.Color) *Diff {
	// check to see if it's legal
	if !b.Legal(start, col) {
		return nil
	}

	// put the stone on the board
	b.Set(start, col)

	// remove dead groups of opposite color
	remove := b.RemoveDead(start, col.Opposite())

	// return diff
	cs := NewCoordSet()
	cs.Add(start)
	add := NewStoneSet(cs, col)
	return NewDiff([]*StoneSet{add}, []*StoneSet{remove})
}

func (b *Board) ApplyDiff(d *Diff) {
	if d == nil {
		return
	}
	for _, add := range d.Add {
		b.SetMany(add.Coords, add.Color)
	}
	for _, remove := range d.Remove {
		b.SetMany(remove.Coords, color.Empty)
	}
}

func (b *Board) CurrentDiff() *Diff {
	black := NewCoordSet()
	white := NewCoordSet()
	for j, row := range b.Points {
		for i, c := range row {
			switch c {
			case color.Black:
				black.Add(NewCoord(i, j))
			case color.White:
				white.Add(NewCoord(i, j))
			}
		}
	}
	addBlack := NewStoneSet(black, color.Black)
	addWhite := NewStoneSet(white, color.White)
	return NewDiff([]*StoneSet{addBlack, addWhite}, nil)
}

type EmptyPointType int

const (
	NotCovered EmptyPointType = iota
	BlackPoint
	WhitePoint
	Dame
)

func (b *Board) FindArea(start *Coord, dead CoordSet) (CoordSet, EmptyPointType) {
	if b.Get(start) != color.Empty {
		return nil, NotCovered
	}

	t := NotCovered

	// initiate the stack
	stack := []*Coord{start}

	// initiate elements
	elts := NewCoordSet()

	// start DFS
	for len(stack) > 0 {
		// pop off the stack
		point := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// add to elements
		elts.Add(point)

		// compute neighbors
		nbs := b.Neighbors(point)
		for _, nb := range nbs {
			if b.Get(nb) == color.Empty && !elts.Has(nb) {
				stack = append(stack, nb)
			} else if b.Get(nb) == color.Black && !dead.Has(nb) {
				t |= BlackPoint
			} else if b.Get(nb) == color.Black {
				t |= WhitePoint
			} else if b.Get(nb) == color.White && !dead.Has(nb) {
				t |= WhitePoint
			} else if b.Get(nb) == color.White {
				t |= BlackPoint
			}
		}
	}
	return elts, t
}

func (b *Board) detectAtariDame(col color.Color, dead, dame CoordSet) CoordSet {
	// first make a copy
	c := b.Copy()

	// then fill in all the dame with a "filler"
	for _, d := range dame.List() {
		c.Set(d, col)
	}

	// these are the "atari" dame to be returned
	points := NewCoordSet()

	// loop
	for {
		changed := false
		// find the (living) groups with 1 liberty
		gps := c.Groups()
		for _, gp := range gps {
			if len(gp.Libs) == 1 {
				rep := gp.Coords.List()[0]
				lib := gp.Libs.List()[0]

				// if it's dead, nothing to do
				if dead.Has(rep) {
					continue
				}
				// add the liberty to the list of atari dame
				points.Add(lib)

				// fill it and keep going
				col := c.Get(rep)
				c.Set(lib, col.Fill())
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return points
}

func (b *Board) DetectAtariDame(dead, dame CoordSet) CoordSet {
	// first detect by filling in dame with FillBlack
	bDame := b.detectAtariDame(color.FillBlack, dead, dame)

	// then FillWhite
	wDame := b.detectAtariDame(color.FillWhite, dead, dame)

	// return the intersection
	return bDame.Intersect(wDame)
}

type ScoreResult struct {
	BlackArea []*Coord
	WhiteArea []*Coord
	BlackDead []*Coord
	WhiteDead []*Coord
	Dame      []*Coord
}

func (b *Board) Score(dead CoordSet, markedDame CoordSet) *ScoreResult {
	blackArea := NewCoordSet()
	whiteArea := NewCoordSet()
	blackDead := NewCoordSet()
	whiteDead := NewCoordSet()
	dame := NewCoordSet()

	// make a new grid to keep track of territory
	grid := [][]EmptyPointType{}

	for i := 0; i < b.Size; i++ {
		grid = append(grid, make([]EmptyPointType, b.Size))
	}

	// add dead stones to the grid, then double count for both area and caps
	for _, coord := range dead {
		switch col := b.Get(coord); col {
		case color.Black:
			grid[coord.Y][coord.X] = WhitePoint
			whiteArea.Add(coord)
			blackDead.Add(coord)
		case color.White:
			grid[coord.Y][coord.X] = BlackPoint
			blackArea.Add(coord)
			whiteDead.Add(coord)
		}
	}

	// add marked dame to dame
	dame.AddAll(markedDame)

	// go through every empty point in the grid
	// anything that hasn't been handled yet gets assigned to either
	// - black area
	// - white area
	// - dame
	for j, row := range b.Points {
		for i, c := range row {
			switch c {
			case color.Black, color.White:
			case color.Empty:
				if grid[j][i] == NotCovered {
					area, t := b.FindArea(NewCoord(i, j), dead)
					for _, coord := range area {
						// this only happens because of marked dame
						if dame.Has(coord) {
							continue
						}
						grid[coord.Y][coord.X] = t
						switch t {
						case BlackPoint:
							blackArea.Add(coord)
						case WhitePoint:
							whiteArea.Add(coord)
						case Dame:
							dame.Add(coord)
						}
					}
				}
			}
		}
	}

	// remove points that need to be filled
	fillAtari := b.DetectAtariDame(dead, dame)
	for _, a := range fillAtari.List() {
		blackArea.Remove(a)
		whiteArea.Remove(a)
	}

	return &ScoreResult{
		BlackArea: blackArea.List(),
		WhiteArea: whiteArea.List(),
		BlackDead: blackDead.List(),
		WhiteDead: whiteDead.List(),
		Dame:      dame.List(),
	}
}
