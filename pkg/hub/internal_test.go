/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/config"
	"github.com/jarednogo/board/pkg/loader"
	"github.com/jarednogo/board/pkg/logx"
	"github.com/jarednogo/board/pkg/twitch"
)

func TestTwitchCallbackSubscribe(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := NewHubWithDB(loader.NewMemoryLoader(), config.Test(), logger)
	assert.NoError(t, err)

	h.tc.SetHTTPClient(twitch.NewMockHTTPClient([]string{
		`{"access_token": "my_user_access_token"}`,
		`{"data": [{"id": "123456789", "login": "some_login"}]}`,
		`{"access_token": "my_app_access_token"}`,
		`{"data": [{"id": "123456789"}]}`,
	}))

	rec := httptest.NewRecorder()

	r := h.TwitchRouter()

	req := httptest.NewRequest(
		"GET",
		"/callback?state=state123&code=anycode&scope=anyscope",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "state123",
	})

	r.ServeHTTP(rec, req)
	assert.Equal(t, rec.Body.String(), "success")
}

func TestTwitchCallbackUnsubscribe(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := NewHubWithDB(loader.NewMemoryLoader(), config.Test(), logger)
	assert.NoError(t, err)

	h.tc.SetHTTPClient(twitch.NewMockHTTPClient([]string{
		`{"access_token": "my_user_access_token"}`,
		`{"data": [{"id": "user123456789", "login": "abc123"}]}`,
		`{"access_token": "my_app_access_token1"}`,
		`{"access_token": "my_app_access_token2"}`,
		`{"total":1, "data": [{"id": "subscription123", "condition": {"broadcaster_user_id": "user123456789"}}]}`,
		`{"access_token": "foobar"}`,
	}))

	rec := httptest.NewRecorder()

	r := h.TwitchRouter()

	req := httptest.NewRequest(
		"GET",
		"/callback?state=state123&code=anycode",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "state123",
	})

	r.ServeHTTP(rec, req)
	assert.Equal(t, rec.Body.String(), "success")
}

func TestTwitchCallbackPost1(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := NewHubWithDB(loader.NewMemoryLoader(), config.Test(), logger)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	r := h.TwitchRouter()

	body := []byte(`{"subscription": {"id": "sub123"}, "event": {"broadcaster_user_id": "broadcaster123", "chatter_user_id": "broadcaster123", "message": {"text": "!setboard abc"}}}`)
	req := httptest.NewRequest("POST", "/callback", bytes.NewBuffer(body))
	r.ServeHTTP(rec, req)

	roomID := h.db.TwitchGetRoom("broadcaster123")
	assert.Equal(t, roomID, "abc")
}

func TestTwitchCallbackPost2(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	h, err := NewHubWithDB(loader.NewMemoryLoader(), config.Test(), logger)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	r := h.TwitchRouter()

	body := []byte(`{"subscription": {"id": "sub123"}, "event": {"broadcaster_user_id": "broadcaster123", "chatter_user_id": "broadcaster123", "message": {"text": "!setboard abc"}}}`)
	req := httptest.NewRequest("POST", "/callback", bytes.NewBuffer(body))
	r.ServeHTTP(rec, req)

	roomID := h.db.TwitchGetRoom("broadcaster123")
	assert.Equal(t, roomID, "abc")

	body = []byte(`{"subscription": {"id": "sub123"}, "event": {"broadcaster_user_id": "broadcaster123", "chatter_user_id": "chatter123", "message": {"text": "!branch e10 e9 e8 e7"}}}`)
	req = httptest.NewRequest("POST", "/callback", bytes.NewBuffer(body))
	r.ServeHTTP(rec, req)

	room, err := h.GetRoom(roomID)
	assert.NoError(t, err)
	save := room.SaveState()
	assert.Equal(t, len(save.SGF), 152)
}
