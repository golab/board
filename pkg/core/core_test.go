/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core_test

import (
	"fmt"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/core"
)

var sanitizeTests = []struct {
	input  string
	output string
}{
	{"  AAA  ", "aaa"},
	{"http://www.example.com?foo=Bar&baz=BOT", "httpwwwexamplecomfoobarbazbot"},
	{"AbC-123", "abc-123"},
	{"MyCoolBoard", "mycoolboard"},
}

func TestSanitize(t *testing.T) {
	for i, tt := range sanitizeTests {
		t.Run(fmt.Sprintf("sanitize%d", i), func(t *testing.T) {
			output := core.Sanitize(tt.input)
			assert.Equal(t, tt.output, output)
		})
	}
}

func TestUUID(t *testing.T) {
	s := core.UUID4()
	assert.Equal(t, len(s), 36)
}
