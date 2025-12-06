/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"fmt"

	"github.com/jarednogo/board/pkg/core/color"
	"github.com/jarednogo/board/pkg/core/coord"
	"github.com/jarednogo/board/pkg/event"
)

func DecodeToCommand(evt event.Event) (Command, error) {
	switch evt.Type() {
	case "add_stone":
		val := evt.Value().(map[string]any)
		c, err := coord.InterfaceToCoord(val["coords"])
		if err != nil {
			return nil, err
		}
		col := color.Color(val["color"].(float64))
		return NewAddStoneCommand(c, col), nil
	case "pass":
		col := color.Color(evt.Value().(float64))
		return NewPassCommand(col), nil
	case "remove_stone":
		c, err := coord.InterfaceToCoord(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewRemoveStoneCommand(c), nil
	case "triangle":
		c, err := coord.InterfaceToCoord(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewAddTriangleCommand(c), nil
	case "square":
		c, err := coord.InterfaceToCoord(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewAddSquareCommand(c), nil
	case "letter":
		val := evt.Value().(map[string]any)
		c, err := coord.InterfaceToCoord(val["coords"])
		if err != nil {
			return nil, err
		}
		letter := val["letter"].(string)
		return NewAddLetterCommand(c, letter), nil
	case "number":
		val := evt.Value().(map[string]any)
		c, err := coord.InterfaceToCoord(val["coords"])
		if err != nil {
			return nil, err
		}
		number := int(val["number"].(float64))
		return NewAddNumberCommand(c, number), nil
	case "remove_mark":
		c, err := coord.InterfaceToCoord(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewRemoveMarkCommand(c), nil
	case "cut":
		return NewCutCommand(), nil
	case "left":
		return NewLeftCommand(), nil
	case "right":
		return NewRightCommand(), nil
	case "up":
		return NewUpCommand(), nil
	case "down":
		return NewDownCommand(), nil
	case "rewind":
		return NewRewindCommand(), nil
	case "fastforward":
		return NewFastForwardCommand(), nil
	case "goto_grid":
		index := int(evt.Value().(float64))
		return NewGotoGridCommand(index), nil
	case "goto_coord":
		c, err := coord.InterfaceToCoord(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewGotoCoordCommand(c), nil
	case "comment":
		val := evt.Value().(string)
		return NewCommentCommand(val), nil
	case "draw":
		vals := evt.Value().([]any)
		var x0 float64
		var y0 float64
		if vals[0] == nil {
			x0 = -1.0
		} else {
			x0 = vals[0].(float64)
		}

		if vals[1] == nil {
			y0 = -1.0
		} else {
			y0 = vals[1].(float64)
		}

		x1 := vals[2].(float64)
		y1 := vals[3].(float64)
		color := vals[4].(string)
		return NewDrawCommand(x0, y0, x1, y1, color), nil
	case "erase_pen":
		return NewErasePenCommand(), nil
	case "copy":
		return NewCopyCommand(), nil
	case "clipboard":
		return NewPasteCommand(), nil
	case "graft":
		// convert the event value to a string and split into tokens
		v := evt.Value().(string)
		return NewGraftCommand(v), nil
	case "score":
		return NewScoreCommand(), nil
	case "markdead":
		c, err := coord.InterfaceToCoord(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewMarkDeadCommand(c), nil
	}

	return nil, fmt.Errorf("unhandled event type")
}
