/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package plugin

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/fetch"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
	"github.com/golab/board/pkg/event"
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
			rank := makeRank(tc.input)
			assert.Equal(t, rank, tc.output)
		})
	}
}

func TestGameinfoToSGF(t *testing.T) {
	data := `{"players": {"black": {"username": "player_black", "rank": 30}, "white": {"username": "player_white", "rank": 29}}, "game_name":"game", "komi":0.5, "width":19, "rules":"chinese"}`

	var payload map[string]any
	err := json.Unmarshal([]byte(data), &payload)
	assert.NoError(t, err)

	o := OGSConnector{}
	sgf := o.gameInfoToSGF(payload)
	_ = sgf
}

func TestGamedata(t *testing.T) {
	data := `{"players": {"black": {"username": "player_black", "rank": 30}, "white": {"username": "player_white", "rank": 29}}, "game_name":"game", "komi":0.5, "width":19, "rules":"chinese", "initial_player":"black", "moves":[[15,15,3150],[3,15,1132],[3,3,1643],[15,3,1150]], "initial_state":{"black":"","white":""}}`

	var payload map[string]any
	err := json.Unmarshal([]byte(data), &payload)
	assert.NoError(t, err)

	o := &OGSConnector{}
	sgf := o.gamedataToSGF(payload)
	assert.Equal(t, sgf, "(;GM[1]FF[4]CA[UTF-8]SZ[19]PB[player_black]PW[player_white]BR[1d]WR[2k]RU[chinese]KM[0.500000]GN[game];B[pp];W[dp];B[dd];W[pd])")
}

func TestOGSConnector1(t *testing.T) {
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r := &MockRoom{}
	_, err := NewMockOGSConnector(r, f)
	assert.NoError(t, err)
}

func TestOGSConnectorSend(t *testing.T) {
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r := &MockRoom{}
	o, err := NewMockOGSConnector(r, f)
	assert.NoError(t, err)
	err = o.send("topic", make(map[string]any))
	assert.NoError(t, err)
}

func TestOGSConnectorConnect1(t *testing.T) {
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r := &MockRoom{}
	o, err := NewMockOGSConnector(r, f)
	assert.NoError(t, err)
	err = o.connect(123456789, "game")
	assert.NoError(t, err)
}

func TestOGSConnectorConnect2(t *testing.T) {
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r := &MockRoom{}
	o, err := NewMockOGSConnector(r, f)
	assert.NoError(t, err)
	err = o.connect(123456789, "review")
	assert.NoError(t, err)
}

func TestOGSConnectorChatConnect(t *testing.T) {
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r := &MockRoom{}
	o, err := NewMockOGSConnector(r, f)
	assert.NoError(t, err)
	err = o.chatConnect()
	assert.NoError(t, err)
}

func TestOGSConnectorReadSocket(t *testing.T) {
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r := &MockRoom{}
	o, err := NewMockOGSConnector(r, f)
	assert.NoError(t, err)

	s := "foobar"
	_, err = o.Socket.Write([]byte(s))
	assert.NoError(t, err)

	socketchan := make(chan byte, len(s))
	err = o.readSocketToChan(socketchan)
	assert.NotNil(t, err)

	data := []byte{}
	for {
		b, ok := <-socketchan
		if !ok {
			break
		}
		data = append(data, b)
	}
	assert.Equal(t, string(data), s)
}

func TestOGSConnectorReadFrame(t *testing.T) {
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r := &MockRoom{}
	o, err := NewMockOGSConnector(r, f)
	assert.NoError(t, err)

	s := "[abc [def] [ghi] [jkl]]"
	_, err = o.Socket.Write([]byte(s))
	assert.NoError(t, err)

	socketchan := make(chan byte, len(s))
	err = o.readSocketToChan(socketchan)
	// should be EOF
	assert.NotNil(t, err)

	data, err := readFrameFromChan(socketchan)
	assert.NoError(t, err)

	assert.Equal(t, string(data), s)
}

func TestOGSConnectorLoop1(t *testing.T) {
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r := &MockRoom{}
	o, err := NewMockOGSConnector(r, f)
	assert.NoError(t, err)

	// setup
	s := `["game/123456789/move", {"move": [3.0, 3.0]}]`
	_, err = o.Socket.Write([]byte(s))
	assert.NoError(t, err)

	// test loop
	err = o.loop(123456789, "game")
	assert.NoError(t, err)

	// check effects
	require.Equal(t, len(r.calls), 3)
	assert.Equal(t, r.calls[0].name, "HeadColor")
	assert.Equal(t, r.calls[1].name, "PushHead")
	assert.Equal(t, r.calls[2].name, "BroadcastFullFrame")
}

func TestOGSConnectorLoop2(t *testing.T) {
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r := &MockRoom{}
	o, err := NewMockOGSConnector(r, f)
	assert.NoError(t, err)

	// setup
	s := `["game/123456789/gamedata", {"width": 19, "komi": 6.5, "game_name": "somename", "rules": "Japanese", "players": {"black": {"username": "playerblack", "rank": 30}, "white": {"username": "playerwhite", "rank": 29}}, "initial_player": "black", "initial_state": {"black": "", "white": ""}, "moves": []}]`
	_, err = o.Socket.Write([]byte(s))
	assert.NoError(t, err)

	// test loop
	err = o.loop(123456789, "game")
	assert.NoError(t, err)

	// check effects
	require.Equal(t, len(r.calls), 2)

	assert.Equal(t, r.calls[0].name, "UploadSGF")
	require.Equal(t, len(r.calls[0].args), 1)
	assert.Equal(t, len(r.calls[0].args[0].(string)), 106)

	assert.Equal(t, r.calls[1].name, "Broadcast")
	require.Equal(t, len(r.calls[1].args), 1)
	evt := r.calls[1].args[0].(event.Event)
	assert.Equal(t, evt.Type(), "testevent")
}

func TestOGSConnectorLoop3(t *testing.T) {
	f := fetch.NewMockFetcher(`{"user": {"id": 123456, "username": "someuser"}, "user_jwt": "jwt123"}`)
	r := &MockRoom{}
	o, err := NewMockOGSConnector(r, f)
	assert.NoError(t, err)

	// setup
	s := `["review/123456789/r", {"f": 2, "m": "aabbccdd"}]`
	_, err = o.Socket.Write([]byte(s))
	assert.NoError(t, err)

	// test loop
	err = o.loop(123456789, "game")
	assert.NoError(t, err)

	// check effects
	require.Equal(t, len(r.calls), 3)
	assert.Equal(t, r.calls[0].name, "GetColorAt")

	assert.Equal(t, r.calls[1].name, "AddStonesToTrunk")
	require.Equal(t, len(r.calls[1].args), 2)
	stones, ok := r.calls[1].args[1].([]*coord.Stone)
	require.True(t, ok)

	require.Equal(t, len(stones), 4)

	assert.Equal(t, stones[0].Color, color.Black)
	assert.True(t, stones[0].Coord.Equal(coord.NewCoord(0, 0)))

	assert.Equal(t, stones[1].Color, color.White)
	assert.True(t, stones[1].Coord.Equal(coord.NewCoord(1, 1)))

	assert.Equal(t, stones[2].Color, color.Black)
	assert.True(t, stones[2].Coord.Equal(coord.NewCoord(2, 2)))

	assert.Equal(t, stones[3].Color, color.White)
	assert.True(t, stones[3].Coord.Equal(coord.NewCoord(3, 3)))

	assert.Equal(t, r.calls[2].name, "BroadcastFullFrame")
}
