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
	"testing"

	"github.com/jarednogo/board/internal/assert"
)

func TestChoices(t *testing.T) {
	numAnimals := len(animals)
	numAdjectives := len(adjectives)
	numColors := len(colors)
	num := numAnimals * numAdjectives * numColors
	// numAnimals == 374
	// numAdjectives == 828
	// numColors == 109
	// 374 * 828 * 109 = 33754248
	// 2^25 = 33554432
	// hopefully 33 million is plenty
	assert.True(t, num > 1<<25)
}

func TestAnimalsLower(t *testing.T) {
	for _, animal := range animals {
		assert.Equal(t, animal, strings.ToLower(animal))
	}
}

func TestAnimalsSplit(t *testing.T) {
	for _, animal := range animals {
		assert.Equal(t, len(strings.Split(animal, " ")), 1)
	}
}

func TestAdjectivesLower(t *testing.T) {
	for _, adjective := range adjectives {
		assert.Equal(t, adjective, strings.ToLower(adjective))
	}
}

func TestAdjectivesSplit(t *testing.T) {
	for _, adjective := range adjectives {
		assert.Equal(t, len(strings.Split(adjective, " ")), 1)
	}
}

func TestColorsLower(t *testing.T) {
	for _, color := range colors {
		assert.Equal(t, color, strings.ToLower(color))
	}
}

func TestColorsSplit(t *testing.T) {
	for _, color := range colors {
		assert.Equal(t, len(strings.Split(color, " ")), 1)
	}
}
