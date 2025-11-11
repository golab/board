/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	main "github.com/jarednogo/board/cmd"
	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/config"
)

func TestPing(t *testing.T) {
	_, r, err := main.Setup(config.Test())
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/ping", nil)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)

	pong := struct {
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(body, &pong)
	assert.NoError(t, err)
	assert.Equal(t, pong.Message, "pong")
}

func TestTwitch(t *testing.T) {
	h, r, err := main.Setup(config.Test())
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	body := []byte(`{"event": {"message": {"text": "!setboard Board"}}}`)
	req := httptest.NewRequest("POST", "/apps/twitch/callback", bytes.NewBuffer(body))
	r.ServeHTTP(rec, req)

	body = []byte(`{"event": {"message": {"text": "!branch k10 k11"}}}`)
	req = httptest.NewRequest("POST", "/apps/twitch/callback", bytes.NewBuffer(body))
	r.ServeHTTP(rec, req)

	room, err := h.GetRoom("Board")
	assert.NoError(t, err)
	assert.Equal(t, len(room.ToSGF()), 77)
}
