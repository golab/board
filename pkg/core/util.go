/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// the core package provides basic functionality to all the major components of the code
package core

import (
	"strings"

	"github.com/google/uuid"
)

// Sanitize ensures strings only contain letters and numbers
func Sanitize(s string) string {
	sb := strings.Builder{}
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			sb.WriteRune(c + 32)
		}
		if (c >= '0' && c <= '9') ||
			(c >= 'a' && c <= 'z') || c == '-' {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

// UUID4 makes and sanitizes a new uuid
func UUID4() string {
	r, _ := uuid.NewRandom()
	s := r.String()
	// remove hyphens
	return Sanitize(s)
}
