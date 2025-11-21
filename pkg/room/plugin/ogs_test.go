/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package plugin_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/room/plugin"
)

func TestMakeRank(t *testing.T) {
	testcases := []struct {
		input  float64
		output string
	}{
		{25, "6k"},
		{31, "2d"},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("MakeRank%d", i), func(t *testing.T) {
			rank := plugin.MakeRank(tc.input)
			assert.Equal(t, rank, tc.output)
		})
	}
}

func TestGameinfoToSGF(t *testing.T) {
	data := `{"players": {"black": {"username": "player_black", "rank": 30}, "white": {"username": "player_white", "rank": 29}}, "game_name":"game", "komi":0.5, "width":19, "rules":"chinese"}`

	var payload map[string]any
	err := json.Unmarshal([]byte(data), &payload)
	assert.NoError(t, err)

	o := plugin.OGSConnector{}
	sgf := o.GameInfoToSGF(payload)
	_ = sgf
}

func TestGamedata(t *testing.T) {
	data := `{"players": {"black": {"username": "player_black", "rank": 30}, "white": {"username": "player_white", "rank": 29}}, "game_name":"game", "komi":0.5, "width":19, "rules":"chinese", "initial_player":"black", "moves":[[15,15,3150],[3,15,1132],[3,3,1643],[15,3,1150]], "initial_state":{"black":"","white":""}}`

	var payload map[string]any
	err := json.Unmarshal([]byte(data), &payload)
	assert.NoError(t, err)

	o := &plugin.OGSConnector{}
	sgf := o.GamedataToSGF(payload)
	assert.Equal(t, sgf, "(;GM[1]FF[4]CA[UTF-8]SZ[19]PB[player_black]PW[player_white]BR[1d]WR[2k]RU[chinese]KM[0.500000]GN[game];B[pp];W[dp];B[dd];W[pd])")
}
