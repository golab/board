/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package logx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/logx"
)

func TestLoggerDebug(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelDebug)
	logger.Debug("foobar")

	var m map[string]any
	err := json.Unmarshal(logger.Bytes(), &m)
	assert.NoError(t, err)
	assert.Equal(t, m["level"], "DEBUG")
	assert.Equal(t, m["msg"], "foobar")
}

func TestLoggerDebugLines(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelDebug)
	logger.Debug("foobar")
	logger.Debug("bazbot")

	lines := logger.Lines()
	assert.Equal(t, len(lines), 2)
}

func TestLoggerInfo(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelInfo)
	logger.Info("foobar")

	var m map[string]any
	err := json.Unmarshal(logger.Bytes(), &m)
	assert.NoError(t, err)
	assert.Equal(t, m["level"], "INFO")
	assert.Equal(t, m["msg"], "foobar")
}

func TestLoggerWarn(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelWarn)
	logger.Warn("foobar")

	var m map[string]any
	err := json.Unmarshal(logger.Bytes(), &m)
	assert.NoError(t, err)
	assert.Equal(t, m["level"], "WARN")
	assert.Equal(t, m["msg"], "foobar")
}

func TestLoggerError(t *testing.T) {
	logger := logx.NewRecorder(logx.LogLevelError)
	logger.Error("foobar")

	var m map[string]any
	err := json.Unmarshal(logger.Bytes(), &m)
	assert.NoError(t, err)
	assert.Equal(t, m["level"], "ERROR")
	assert.Equal(t, m["msg"], "foobar")
}

func TestMiddleware(t *testing.T) {
	// Create logger at Debug level so Debug() messages are emitted.
	logger := logx.NewRecorder(logx.LogLevelDebug)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
		_, _ = w.Write([]byte("ok"))
	})

	mw := logger.AsMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	rec := httptest.NewRecorder()

	mw.ServeHTTP(rec, req)

	assert.True(t, called)

	// find the first line that has msg == "http"
	var m map[string]any
	err := json.Unmarshal(logger.Bytes(), &m)
	assert.NoError(t, err)
	assert.Equal(t, m["level"], "DEBUG")
	assert.Equal(t, m["msg"], "http")
	assert.Equal(t, m["method"], "GET")
	assert.Equal(t, m["path"], "/foo")
}
