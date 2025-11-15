/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package logx

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type DefaultLogger struct {
	s *slog.Logger
}

func NewDefaultLogger() *DefaultLogger {
	handler := slog.NewJSONHandler(os.Stderr, nil)
	return &DefaultLogger{
		slog.New(handler),
	}
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
