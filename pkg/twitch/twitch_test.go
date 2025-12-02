/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package twitch_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/twitch"
)

func MockHTTPResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

type MockHTTPClient struct {
	responses []string
	index     int
}

func NewMockHTTPClient(responses []string) *MockHTTPClient {
	return &MockHTTPClient{responses: responses}
}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	response := ""
	if c.index < len(c.responses) {
		response = c.responses[c.index]
		c.index++
	}
	return MockHTTPResponse(200, response), nil
}

func TestGetUserAccessToken(t *testing.T) {
	tc := twitch.NewDefaultTwitchClient("", "", "", "")
	responses := []string{`{"access_token": "foobar"}`}
	tc.SetHTTPClient(NewMockHTTPClient(responses))
	token, err := tc.GetUserAccessToken("")
	assert.NoError(t, err)
	assert.Equal(t, token, "foobar")
}

func TestGetUsers(t *testing.T) {
	tc := twitch.NewDefaultTwitchClient("", "", "", "")
	responses := []string{`{"data": [{"id": "123456789", "login": "some_login"}]}`}
	tc.SetHTTPClient(NewMockHTTPClient(responses))
	user, err := tc.GetUsers("")
	assert.NoError(t, err)
	assert.Equal(t, user, "123456789")
}

func TestGetAppAccessToken(t *testing.T) {
	tc := twitch.NewDefaultTwitchClient("", "", "", "")
	responses := []string{`{"access_token": "foobar"}`}
	tc.SetHTTPClient(NewMockHTTPClient(responses))
	user, err := tc.GetAppAccessToken()
	assert.NoError(t, err)
	assert.Equal(t, user, "foobar")
}

func TestUnsubscribe(t *testing.T) {
	tc := twitch.NewDefaultTwitchClient("", "", "", "")
	responses := []string{`{"access_token": "foobar"}`}
	tc.SetHTTPClient(NewMockHTTPClient(responses))
	err := tc.Unsubscribe("", "")
	assert.NoError(t, err)
}

func TestSubscribe(t *testing.T) {
	tc := twitch.NewDefaultTwitchClient("", "", "", "")
	responses := []string{`{"data": [{"id": "123456789"}]}`}
	tc.SetHTTPClient(NewMockHTTPClient(responses))
	id, err := tc.Subscribe("", "")
	assert.NoError(t, err)
	assert.Equal(t, id, "123456789")
}

func TestAlreadySubscribed(t *testing.T) {
	tc := twitch.NewDefaultTwitchClient("", "", "", "")
	responses := []string{`{"message": "subscription exists"}`}
	tc.SetHTTPClient(NewMockHTTPClient(responses))
	_, err := tc.Subscribe("", "")
	assert.ErrorIs(t, err, fmt.Errorf("subscription exists"))
}

func TestGetSubscription(t *testing.T) {
	tc := twitch.NewDefaultTwitchClient("", "", "", "")
	responses := []string{
		`{"access_token": "foobar"}`,
		`{"total":1, "data": [{"id": "subscription123", "condition": {"broadcaster_user_id": "abc123"}}]}`,
	}
	tc.SetHTTPClient(NewMockHTTPClient(responses))
	id, err := tc.GetSubscription("abc123")
	assert.NoError(t, err)
	assert.Equal(t, id, "subscription123")
}

func TestVerify(t *testing.T) {
	secret := "9piitrv7ch5yyr56b0cbct5t9bli92"
	message := "hello world"
	sig := "sha256=3ae321e96e012c6cc89b73326f329f0a3d1d7935abc3d819387830ce5f1b3074"
	assert.True(t, twitch.Verify(secret, message, sig))
}

func TestParse(t *testing.T) {
	testcases := []struct {
		input    string
		haserror bool
		command  string
		body     string
	}{
		{"!a b c d", false, "a", "b c d"},
		{"", true, "", ""},
		{"cmd", true, "", ""},
		{"!", true, "", ""},
		{"!foo bar", false, "foo", "bar"},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("twitchParse%d", i), func(t *testing.T) {
			chat, err := twitch.Parse(tc.input)
			if err != nil {
				assert.True(t, tc.haserror)
				return
			}
			assert.Equal(t, chat.Command, tc.command)
			assert.Equal(t, chat.Body, tc.body)
		})
	}
}

func FuzzParseChat(f *testing.F) {
	testcases := []string{"", "!foo bar", "!", "foo bar"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig string) {
		// looking for crashes or panics
		_, _ = twitch.Parse(orig)
	})
}
