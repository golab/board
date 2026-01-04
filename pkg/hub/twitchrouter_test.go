/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/pkg/config"
	"github.com/golab/board/pkg/hub"
	"github.com/golab/board/pkg/logx"
)

func TestTwitchRouterCallbackPost(t *testing.T) {
	twitchTests := []struct {
		body    string
		code    int
		numLogs int
		output  string
	}{
		{"foobar", 200, 1, ""},
		{`{"challenge": "foo"}`, 200, 1, "foo"},
		{`{"abc": "123"}`, 200, 1, ""},
		{`{"event": {}}`, 200, 1, ""},
		{`{"event": {"message": {"text": "abc"}}}`, 200, 1, ""},
		{`{"event": {"message": {"text": "!abc"}}}`, 200, 1, ""},
		{`{"event": {"message": {"text": "!setboard"}}}`, 200, 2, ""},
		{`{"event": {"broadcaster_user_id": "x", "chatter_user_id": "y", "message": {"text": "!setboard abc"}}}`, 200, 2, ""},
		{`{"event": {"broadcaster_user_id": "x", "chatter_user_id": "y", "message": {"text": "!branch a1 a2"}}}`, 200, 2, ""},
	}
	for i, tt := range twitchTests {
		t.Run(fmt.Sprintf("twitch%d", i), func(t *testing.T) {

			logger := logx.NewRecorder(logx.LogLevelDebug)
			h, err := hub.NewHub(config.Test(), logger)
			require.NoError(t, err)
			r := chi.NewRouter()
			r.Mount("/", h.TwitchRouter())

			postBody := []byte(tt.body)
			req := httptest.NewRequest("POST", "/callback", bytes.NewBuffer(postBody))
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, len(logger.Lines()), tt.numLogs)

			body, err := io.ReadAll(rec.Body)
			require.NoError(t, err)
			assert.Equal(t, string(body), tt.output)
		})
	}
}

func TestTwitchRouterCallbackGet(t *testing.T) {
	twitchTests := []struct {
		code    int
		args    string
		numLogs int
		output  string
	}{
		{200, "", 1, "invalid state\n"},
		{200, "?state=foo", 1, "success"},
	}
	for i, tt := range twitchTests {
		t.Run(fmt.Sprintf("twitch%d", i), func(t *testing.T) {

			logger := logx.NewRecorder(logx.LogLevelDebug)
			h, err := hub.NewHub(config.Test(), logger)
			require.NoError(t, err)
			r := chi.NewRouter()
			r.Mount("/", h.TwitchRouter())

			req := httptest.NewRequest("GET", "/callback"+tt.args, nil)
			req.Header.Add("Cookie", "oauth_state=foo")
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, len(logger.Lines()), tt.numLogs)
			t.Logf("%v", logger.Lines())

			body, err := io.ReadAll(rec.Body)
			require.NoError(t, err)
			assert.Equal(t, string(body), tt.output)
			t.Logf("%v", string(body))
		})
	}
}
