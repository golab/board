/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package plugin_test

import (
	"encoding/json"
	"testing"
)

func TestGamedata(t *testing.T) {
	data := `{"white_player_id":0, "black_player_id":1, "game_name":"game", "komi":0.5, "width":19, "rules":"chinese", "initial_player":"black", "moves":[[15,15,3150],[3,15,1132],[3,3,1643],[15,3,1150]], "initial_state":{"black":"","white":""}}`

	var payload any
	err := json.Unmarshal([]byte(data), &payload)
	if err != nil {
		t.Error(err)
	}
	_, ok := payload.(map[string]any)
	if !ok {
		t.Errorf("error while coercing interface to map[string]any")
	}

	// currently no way to do tests without making network connections
}
