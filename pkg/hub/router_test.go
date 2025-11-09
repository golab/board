/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/config"
	"github.com/jarednogo/board/pkg/hub"
)

func TestApiRouter(t *testing.T) {
	h, err := hub.NewHub(config.Test())
	assert.NoError(t, err, "apirouter")
	_ = h
	r := chi.NewRouter()
	r.Mount("/api", hub.ApiRouter("version"))

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/version")
	if err != nil {
		assert.NoError(t, err, "version")
		return
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)
	pong := struct {
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(body, &pong)
	assert.NoError(t, err, "unmarshal")

	if pong.Message != "version" {
		t.Fatalf("expected version, got %s", body)
	}
}
