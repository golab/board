/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package integration_test

import (
	"testing"

	"github.com/golab/board/integration"
	"github.com/golab/board/internal/assert"
	"github.com/golab/board/pkg/core/parser"
	"github.com/golab/board/pkg/state"
)

func TestGenerate(t *testing.T) {
	seed := 0
	num := 500
	minLength := 100
	maxLength := 300
	sgfs := integration.GenerateNRandomSGF(seed, num, minLength, maxLength)

	sgf := parser.Merge(sgfs)
	s, err := state.FromSGF(sgf)
	assert.NoError(t, err)
	//assert.True(t, duration < 2*time.Second)
	// this computation takes a while, so omitting it for now
	//assert.Equal(t, len(s.ToSGF()), 616921)
	assert.Equal(t, len(s.Nodes()), 99977)
}

func sgfsForBenchmark(seed, num int) []string {
	minLength := 100
	maxLength := 300
	return integration.GenerateNRandomSGF(seed, num, minLength, maxLength)
}

// benchmarking how long it takes to turn one sgf into a state
func BenchmarkFromSGF(b *testing.B) {
	sgfs := sgfsForBenchmark(0, 500)
	i := 0
	for b.Loop() {
		sgf := sgfs[i%len(sgfs)]
		state.FromSGF(sgf) //nolint:errcheck
		i++
	}
}

// benchmarking how long it takes to turn one state into an sgf
func BenchmarkToSGF(b *testing.B) {
	sgfs := sgfsForBenchmark(0, 500)
	states := []*state.State{}
	for _, sgf := range sgfs {
		s, err := state.FromSGF(sgf)
		if err != nil {
			b.Error(err)
		}
		states = append(states, s)
	}
	i := 0
	for b.Loop() {
		s := states[i%len(states)]
		s.ToSGF()
		i++
	}
}

// benchmarking how long it takes to merge 100 sgfs
func BenchmarkMerge(b *testing.B) {
	mergeList := [][]string{}
	for i := 0; i < 10; i++ {
		sgfs := sgfsForBenchmark(i, 100)
		mergeList = append(mergeList, sgfs)
	}
	i := 0
	for b.Loop() {
		sgfs := mergeList[i%len(mergeList)]
		parser.Merge(sgfs)
		i++
	}
}

// benchmarking how long it takes to turn a 100-sgf merge into a state
func BenchmarkMergeFromSGF(b *testing.B) {
	mergedSGFs := []string{}
	for i := 0; i < 10; i++ {
		sgfs := sgfsForBenchmark(i, 100)
		sgf := parser.Merge(sgfs)
		mergedSGFs = append(mergedSGFs, sgf)
	}
	i := 0
	for b.Loop() {
		sgf := mergedSGFs[i%len(mergedSGFs)]
		state.FromSGF(sgf) //nolint:errcheck
		i++
	}
}

// benchmarking how long it takes to turn a state with 100 merged sgfs into an sgf
func BenchmarkMergeToSGF(b *testing.B) {
	states := []*state.State{}
	for i := 0; i < 10; i++ {
		sgfs := sgfsForBenchmark(i, 100)
		sgf := parser.Merge(sgfs)
		s, err := state.FromSGF(sgf)
		if err != nil {
			b.Error(err)
		}
		states = append(states, s)
	}

	i := 0
	for b.Loop() {
		s := states[i%len(states)]
		s.ToSGF()
		i++
	}
}

func BenchmarkParse(b *testing.B) {
	mergedSGFs := []string{}
	for i := 0; i < 10; i++ {
		sgfs := sgfsForBenchmark(i, 100)
		sgf := parser.Merge(sgfs)
		mergedSGFs = append(mergedSGFs, sgf)
	}
	i := 0
	for b.Loop() {
		sgf := mergedSGFs[i%len(mergedSGFs)]
		p := parser.New(sgf)
		p.Parse() //nolint:errcheck
		i++
	}
}
