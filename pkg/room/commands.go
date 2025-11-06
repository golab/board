/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package room

import (
	"fmt"

	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/state"
)

func decodeToCommand(evt *core.EventJSON) (state.Command, error) {
	switch evt.Event {
	case "add_stone":
		val := evt.Value.(map[string]interface{})
		c, err := core.InterfaceToCoord(val["coords"])
		if err != nil {
			return nil, err
		}
		col := core.Color(val["color"].(float64))
		return state.NewAddStoneCommand(c, col), nil
	case "pass":
		col := core.Color(evt.Value.(float64))
		return state.NewPassCommand(col), nil
	case "remove_stone":
		c, err := core.InterfaceToCoord(evt.Value)
		if err != nil {
			return nil, err
		}
		return state.NewRemoveStoneCommand(c), nil
	case "triangle":
		c, err := core.InterfaceToCoord(evt.Value)
		if err != nil {
			return nil, err
		}
		return state.NewAddTriangleCommand(c), nil
	case "square":
		c, err := core.InterfaceToCoord(evt.Value)
		if err != nil {
			return nil, err
		}
		return state.NewAddSquareCommand(c), nil
	case "letter":
		val := evt.Value.(map[string]interface{})
		c, err := core.InterfaceToCoord(val["coords"])
		if err != nil {
			return nil, err
		}
		letter := val["letter"].(string)
		return state.NewAddLetterCommand(c, letter), nil
	case "number":
		val := evt.Value.(map[string]interface{})
		c, err := core.InterfaceToCoord(val["coords"])
		if err != nil {
			return nil, err
		}
		number := int(val["number"].(float64))
		return state.NewAddNumberCommand(c, number), nil
	case "remove_mark":
		c, err := core.InterfaceToCoord(evt.Value)
		if err != nil {
			return nil, err
		}
		return state.NewRemoveMarkCommand(c), nil
	case "cut":
		return state.NewCutCommand(), nil
	case "left":
		return state.NewLeftCommand(), nil
	case "right":
		return state.NewRightCommand(), nil
	case "up":
		return state.NewUpCommand(), nil
	case "down":
		return state.NewDownCommand(), nil
	case "rewind":
		return state.NewRewindCommand(), nil
	case "fastforward":
		return state.NewFastForwardCommand(), nil
	case "goto_grid":
		index := int(evt.Value.(float64))
		return state.NewGotoGridCommand(index), nil
	case "goto_coord":
		c, err := core.InterfaceToCoord(evt.Value)
		if err != nil {
			return nil, err
		}
		return state.NewGotoCoordCommand(c), nil
	case "comment":
		val := evt.Value.(string)
		return state.NewCommentCommand(val), nil
	case "draw":
		vals := evt.Value.([]interface{})
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
		return state.NewDrawCommand(x0, y0, x1, y1, color), nil
	case "erase_pen":
		return state.NewErasePenCommand(), nil
	case "copy":
		return state.NewCopyCommand(), nil
	case "clipboard":
		return state.NewClipboardCommand(), nil
	case "graft":
		// convert the event value to a string and split into tokens
		v := evt.Value.(string)
		return state.NewGraftCommand(v), nil
	case "score":
		return state.NewScoreCommand(), nil
	case "markdead":
		c, err := core.InterfaceToCoord(evt.Value)
		if err != nil {
			return nil, err
		}
		return state.NewMarkDeadCommand(c), nil
	}

	return nil, fmt.Errorf("unhandled event type")
}
