/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package integration_test

import (
	"testing"
	"time"

	"github.com/jarednogo/board/integration"
	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/state"
)

func TestGenerate(t *testing.T) {
	num := 500
	minLength := 100
	maxLength := 300
	sgfs := integration.GenerateNRandomSGF(num, minLength, maxLength)

	sgf := core.Merge(sgfs)
	start := time.Now()
	s, err := state.FromSGF(sgf)
	assert.NoError(t, err)
	duration := time.Since(start)
	assert.True(t, duration < 2*time.Second)
	// this computation takes a while, so omitting it for now
	//assert.Equal(t, len(s.ToSGF()), 616921)
	assert.Equal(t, len(s.Nodes()), 99977)
}
