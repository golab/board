/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package logx

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
	With(string, string) Logger
	AsMiddleware(http.Handler) http.Handler
}

type DefaultLogger struct {
	s *slog.Logger
}

func NewDefaultLoggerWithWriter(lvl LogLevel, w io.Writer) *DefaultLogger {
	level := new(slog.LevelVar)
	switch lvl {
	case LogLevelDebug:
		level.Set(slog.LevelDebug)
	case LogLevelInfo:
		level.Set(slog.LevelInfo)
	case LogLevelWarn:
		level.Set(slog.LevelWarn)
	case LogLevelError:
		level.Set(slog.LevelError)
	}

	handler := slog.NewJSONHandler(
		w,
		&slog.HandlerOptions{
			Level: level,
		},
	)
	return &DefaultLogger{
		slog.New(handler),
	}
}

func NewDefaultLogger(lvl LogLevel) *DefaultLogger {
	return NewDefaultLoggerWithWriter(lvl, os.Stdout)
}

func (l *DefaultLogger) With(key, value string) Logger {
	return &DefaultLogger{l.s.With(slog.String(key, value))}
}

func (l *DefaultLogger) Debug(msg string, args ...any) {
	l.s.Debug(msg, args...)
}

func (l *DefaultLogger) Info(msg string, args ...any) {
	l.s.Info(msg, args...)
}

func (l *DefaultLogger) Warn(msg string, args ...any) {
	l.s.Warn(msg, args...)
}

func (l *DefaultLogger) Error(msg string, args ...any) {
	l.s.Error(msg, args...)
}

func (l *DefaultLogger) AsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Debug("http", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type Recorder struct {
	*bytes.Buffer
	*DefaultLogger
}

func NewRecorder(lvl LogLevel) *Recorder {
	var buffer bytes.Buffer
	l := NewDefaultLoggerWithWriter(lvl, &buffer)
	return &Recorder{&buffer, l}
}

func (r *Recorder) Lines() []string {
	s := strings.TrimSpace(r.String())
	return strings.Split(s, "\n")
}

func (r *Recorder) With(key, value string) Logger {
	l := &DefaultLogger{r.s.With(slog.String(key, value))}
	return &Recorder{Buffer: nil, DefaultLogger: l}
}
