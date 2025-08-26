package server

import (
	"github.com/jarednogo/board/backend/core"
)

// helper functions for ogs

func (r *Room) HeadColor() core.Color {
	return r.State.Head.Color
}

func (r *Room) PushHead(x, y int, c core.Color) {
	r.State.PushHead(x, y, c)
}

func (r *Room) GenerateFullFrame(t core.TreeJSONType) *core.Frame {
	return r.State.GenerateFullFrame(t)
}

func (r *Room) AddPatternNodes(movesArr []*core.PatternMove) {
	r.State.AddPatternNodes(movesArr)
}
