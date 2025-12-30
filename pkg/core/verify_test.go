/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core_test

import (
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/pkg/core"
)

func TestHash1(t *testing.T) {
	hash := "$2a$10$O5gxBye06fldX9Al9Upc8.nYkE33KWMTfxYF/sMt5TgEs66Vcg0JK"
	password := "deadbeef"

	if !core.CorrectPassword(password, hash) {
		t.Errorf("error checking password")
	}
}

func TestHash2(t *testing.T) {
	h := core.Hash("test")
	assert.Equal(t, len(h), 60)
}
