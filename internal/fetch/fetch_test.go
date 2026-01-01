/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package fetch_test

import (
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/fetch"
	"github.com/golab/board/internal/sgfsamples"
)

var ogsTests = []struct {
	url string
	ogs bool
}{
	{"https://online-go.com/game/00000001", true},
	{"http://online-go.com/game/00000001", true},
	{"https://test.com/game/00000001", false},
}

func TestOGS(t *testing.T) {
	for _, tt := range ogsTests {
		t.Run(tt.url, func(t *testing.T) {
			if fetch.IsOGS(tt.url) != tt.ogs {
				t.Errorf("error checking for ogs url: %s", tt.url)
			}
		})
	}
}

func TestOGSCheckEnded(t *testing.T) {
	mc := fetch.NewMockClient("")
	f := fetch.NewDefaultFetcher(mc)
	ended, err := f.OGSCheckEnded("https://online-go.com/game/123")
	assert.NoError(t, err)
	assert.True(t, ended)
}

func TestFetchOGS(t *testing.T) {
	mc := fetch.NewMockClient(sgfsamples.Resignation1)
	f := fetch.NewDefaultFetcher(mc)
	sgf, err := f.FetchOGS("https://online-go.com/game/123")
	assert.NoError(t, err)
	assert.Equal(t, len(sgf), len(sgfsamples.Resignation1))
}

func TestFetchOGSReview(t *testing.T) {
	mc := fetch.NewMockClient(sgfsamples.Resignation1)
	f := fetch.NewDefaultFetcher(mc)
	sgf, err := f.FetchOGS("https://online-go.com/review/123")
	assert.NoError(t, err)
	assert.Equal(t, len(sgf), len(sgfsamples.Resignation1))
}

func TestApprovedOGSFetch(t *testing.T) {
	mc := fetch.NewMockClient(sgfsamples.Resignation1)
	f := fetch.NewDefaultFetcher(mc)
	sgf, err := f.ApprovedFetch("https://online-go.com/game/123")
	assert.NoError(t, err)
	assert.Equal(t, len(sgf), len(sgfsamples.Resignation1))
}

func TestApprovedFetch(t *testing.T) {
	mc := fetch.NewMockClient(sgfsamples.Resignation1)
	f := fetch.NewDefaultFetcher(mc)
	sgf, err := f.ApprovedFetch("https://files.gokgs.com/foo.sgf")
	assert.NoError(t, err)
	assert.Equal(t, len(sgf), len(sgfsamples.Resignation1))
}

func TestFetchOther(t *testing.T) {
	mc := fetch.NewMockClient(sgfsamples.Resignation1)
	f := fetch.NewDefaultFetcher(mc)
	sgf, err := f.Fetch("https://gokifu.com/somesgf")
	assert.NoError(t, err)
	assert.Equal(t, len(sgf), len(sgfsamples.Resignation1))
}

func TestFetchForbidden(t *testing.T) {
	mc := fetch.NewMockClient(sgfsamples.Resignation1)
	f := fetch.NewDefaultFetcher(mc)
	_, err := f.ApprovedFetch("https://thesgfdatabase.com/somesgf")
	assert.NotNil(t, err)
}

func TestParseError(t *testing.T) {
	mc := fetch.NewMockClient(sgfsamples.Resignation1)
	f := fetch.NewDefaultFetcher(mc)
	_, err := f.Fetch("---abc")
	assert.NotNil(t, err)
}

func TestBadOGS(t *testing.T) {
	mc := fetch.NewMockClient("")
	f := fetch.NewDefaultFetcher(mc)
	_, err := f.ApprovedFetch("https://online-go.com/foo/bar/baz")
	assert.NotNil(t, err)
}
