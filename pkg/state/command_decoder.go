/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"fmt"

	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
	"github.com/golab/board/pkg/event"
)

func DecodeToCommand(evt event.Event) (Command, error) {
	switch evt.Type() {
	case "add_stone":
		val, ok := evt.Value().(map[string]any)
		if !ok {
			return nil, fmt.Errorf("event requires map for 'value'")
		}
		c, err := coord.FromInterface(val["coords"])
		if err != nil {
			return nil, err
		}
		f, ok := val["color"].(float64)
		if !ok {
			return nil, fmt.Errorf("color should be 1 (black) or 2 (white)")
		}
		col := color.Color(f)
		if col != color.Black && col != color.White {
			return nil, fmt.Errorf("invalid color")
		}
		return NewAddStoneCommand(c, col), nil
	case "pass":
		f, ok := evt.Value().(float64)
		if !ok {
			return nil, fmt.Errorf("value should be 1 (black) or 2 (white)")
		}
		col := color.Color(f)
		if col != color.Black && col != color.White {
			return nil, fmt.Errorf("invalid color")
		}

		return NewPassCommand(col), nil
	case "remove_stone":
		c, err := coord.FromInterface(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewRemoveStoneCommand(c), nil
	case "triangle":
		c, err := coord.FromInterface(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewAddTriangleCommand(c), nil
	case "square":
		c, err := coord.FromInterface(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewAddSquareCommand(c), nil
	case "letter":
		val, ok := evt.Value().(map[string]any)
		if !ok {
			return nil, fmt.Errorf("event requires map for 'value'")
		}
		c, err := coord.FromInterface(val["coords"])
		if err != nil {
			return nil, err
		}
		letter, ok := val["letter"].(string)
		if !ok {
			return nil, fmt.Errorf("'letter' should be a string")
		}
		return NewAddLetterCommand(c, letter), nil
	case "number":
		val, ok := evt.Value().(map[string]any)
		if !ok {
			return nil, fmt.Errorf("event requires map for 'value'")
		}
		c, err := coord.FromInterface(val["coords"])
		if err != nil {
			return nil, err
		}
		f, ok := val["number"].(float64)
		if !ok {
			return nil, fmt.Errorf("'number' should be a number")
		}
		number := int(f)
		return NewAddNumberCommand(c, number), nil
	case "label":
		val, ok := evt.Value().(map[string]any)
		if !ok {
			return nil, fmt.Errorf("event requires map for 'value'")
		}

		c, err := coord.FromInterface(val["coords"])
		if err != nil {
			return nil, err
		}
		label, ok := val["label"].(string)
		if !ok {
			return nil, fmt.Errorf("'label' should be a string")
		}
		return NewAddLabelCommand(c, label), nil
	case "remove_mark":
		c, err := coord.FromInterface(evt.Value())
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
		f, ok := evt.Value().(float64)
		if !ok {
			return nil, fmt.Errorf("'value' should be a number")
		}
		index := int(f)
		return NewGotoGridCommand(index), nil
	case "goto_coord":
		c, err := coord.FromInterface(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewGotoCoordCommand(c), nil
	case "comment":
		val, ok := evt.Value().(string)
		if !ok {
			return nil, fmt.Errorf("'value' should be a string")
		}
		return NewCommentCommand(val), nil
	case "draw":
		vals, ok := evt.Value().([]any)
		if !ok {
			return nil, fmt.Errorf("'value' should be in the format [float, float, float, float, string]")
		}
		var x0 float64
		var y0 float64
		if vals[0] == nil {
			x0 = -1.0
		} else {
			x0, ok = vals[0].(float64)
			if !ok {
				return nil, fmt.Errorf("'value' should be in the format [float, float, float, float, string]")
			}
		}

		if vals[1] == nil {
			y0 = -1.0
		} else {
			y0, ok = vals[1].(float64)
			if !ok {
				return nil, fmt.Errorf("'value' should be in the format [float, float, float, float, string]")
			}

		}

		x1, ok := vals[2].(float64)
		if !ok {
			return nil, fmt.Errorf("'value' should be in the format [float, float, float, float, string]")
		}

		y1, ok := vals[3].(float64)
		if !ok {
			return nil, fmt.Errorf("'value' should be in the format [float, float, float, float, string]")
		}

		color, ok := vals[4].(string)

		if !ok {
			return nil, fmt.Errorf("'value' should be in the format [float, float, float, float, string]")
		}
		return NewDrawCommand(x0, y0, x1, y1, color), nil
	case "erase_pen":
		return NewErasePenCommand(), nil
	case "copy":
		return NewCopyCommand(), nil
	case "clipboard":
		return NewPasteCommand(), nil
	case "graft":
		// convert the event value to a string and split into tokens
		v, ok := evt.Value().(string)
		if !ok {
			return nil, fmt.Errorf("'value' should be a string")
		}
		return NewGraftCommand(v), nil
	case "score":
		return NewScoreCommand(), nil
	case "markdead":
		c, err := coord.FromInterface(evt.Value())
		if err != nil {
			return nil, err
		}
		return NewMarkDeadCommand(c), nil
	}

	return nil, fmt.Errorf("unhandled event type")
}
