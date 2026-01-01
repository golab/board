/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func makeHTTPResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

type MockClient struct {
	sgf string
}

func NewMockClient(sgf string) *MockClient {
	return &MockClient{sgf: sgf}
}

func (mc *MockClient) Get(s string) (*http.Response, error) {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return nil, err
	}
	if u.Hostname() == "online-go.com" {
		paths := strings.Split(u.Path, "/")
		if len(paths) == 5 && paths[1] == "api" && paths[2] == "v1" && paths[3] == "games" {
			return makeHTTPResponse(`{"ended": "true", "source": "sgf"}`), nil

		}
		if len(paths) == 6 && paths[1] == "api" && paths[2] == "v1" && paths[5] == "sgf" {
			return makeHTTPResponse(mc.sgf), nil
		}
		return nil, fmt.Errorf("bad request")
	}
	return makeHTTPResponse(mc.sgf), nil
}
