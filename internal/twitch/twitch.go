/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package twitch

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Chat is used to encode strings like "!command plus some args"
type Chat struct {
	Command string
	Body    string
}

// Parse turns a string into a Chat
func Parse(chat string) (*Chat, error) {
	chat = strings.TrimSpace(chat)
	if !strings.HasPrefix(chat, "!") {
		return nil, fmt.Errorf("not a command")
	}
	chat = strings.TrimPrefix(chat, "!")
	tokens := strings.Split(chat, " ")
	if len(tokens) == 0 {
		return nil, fmt.Errorf("bad command")
	}
	command := strings.ToLower(tokens[0])
	body := strings.Join(tokens[1:], " ")
	return &Chat{command, body}, nil
}

// TwitchJSON is the base response type for twitch messages
type TwitchJSON struct {
	Subscription *Subscription `json:"subscription"`
	Event        *TwitchEvent  `json:"event"`
	Challenge    string        `json:"challenge"`
}

// Subscription only is used to get the id
type Subscription struct {
	ID string `json:"id"`
}

// TwitchEvent is the twitch event type
type TwitchEvent struct {
	BroadcasterUserID string             `json:"broadcaster_user_id"`
	ChatterUserID     string             `json:"chatter_user_id"`
	Message           *TwitchMessageJSON `json:"message"`
}

// TwitchMessageJSON is the twitch message type
type TwitchMessageJSON struct {
	Text string `json:"text"`
}

func Verify(secret, message, signature string) bool {
	if len(secret) == 0 {
		return true
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))

	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

type TwitchClient interface {
	GetUserAccessToken(string) (string, error)
	GetUsers(string) (string, error)
	GetAppAccessToken() (string, error)
	Unsubscribe(string, string) error
	Subscribe(string, string) (string, error)
	GetSubscription(string) (string, error)
	SetHTTPClient(HTTPClient)
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type DefaultTwitchClient struct {
	clientID string
	secret   string
	botID    string
	url      string
	client   HTTPClient
}

func NewDefaultTwitchClient(clientID, secret, botID, url string) *DefaultTwitchClient {
	return &DefaultTwitchClient{
		clientID: clientID,
		secret:   secret,
		botID:    botID,
		url:      url,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *DefaultTwitchClient) SetHTTPClient(client HTTPClient) {
	t.client = client
}

// GetUserAccessToken is a prescribed pattern from twitch to get an access token
// attached to a particular user
func (t *DefaultTwitchClient) GetUserAccessToken(code string) (string, error) {
	body := map[string]string{
		"client_id":     t.clientID,
		"client_secret": t.secret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  fmt.Sprintf("%s/apps/twitch/callback", t.url),
		"scope":         "channel:bot",
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	bodyReader := bytes.NewBuffer(jsonData)
	url := "https://id.twitch.tv/oauth2/token"

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //nolint: errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var s struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return "", err
	}

	return s.AccessToken, nil
}

// GetUsers specifically requires a user access token
// TODO: getting the user access token should just be part of this function
func (t *DefaultTwitchClient) GetUsers(token string) (string, error) {
	url := "https://api.twitch.tv/helix/users"

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Client-Id", t.clientID)

	// Send the request
	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //nolint: errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var s map[string][]struct {
		ID    string `json:"id"`
		Login string `json:"login"`
	}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return "", err
	}

	if _, ok := s["data"]; !ok {
		return "", fmt.Errorf("invalid data returned")
	}

	if len(s["data"]) == 0 {
		return "", fmt.Errorf("no users returned")
	}

	return s["data"][0].ID, nil
}

// GetAppAccessToken is a prescribed pattern from twitch to get an access token
// associated with an application
func (t *DefaultTwitchClient) GetAppAccessToken() (string, error) {
	body := map[string]string{
		"client_id":     t.clientID,
		"client_secret": t.secret,
		"grant_type":    "client_credentials",
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	bodyReader := bytes.NewBuffer(jsonData)
	url := "https://id.twitch.tv/oauth2/token"

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //nolint: errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var s struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return "", err
	}

	return s.AccessToken, nil
}

// SubscriptionRequest  is used to subscribe
type SubscriptionRequest struct {
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition map[string]string `json:"condition"`
	Transport map[string]string `json:"transport"`
}

// Unsubscribe from a twitch channel, requires an app access token
// TODO: getting the access token should be part of this function
func (t *DefaultTwitchClient) Unsubscribe(id, token string) error {
	body := map[string]string{"id": id}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewBuffer(jsonData)
	url := "https://api.twitch.tv/helix/eventsub/subscriptions"

	req, err := http.NewRequest(http.MethodDelete, url, bodyReader)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Client-Id", t.clientID)

	// Send the request
	_, err = t.client.Do(req)
	return err
}

// Subscribe to a channel, requires an app access token
// TODO: getting the access token should be part of this function
func (t *DefaultTwitchClient) Subscribe(user, token string) (string, error) {
	body := SubscriptionRequest{
		Type:    "channel.chat.message",
		Version: "1",
		Condition: map[string]string{
			"broadcaster_user_id": user,
			"user_id":             t.botID,
		},
		Transport: map[string]string{
			"method":   "webhook",
			"callback": fmt.Sprintf("%s/apps/twitch/callback", t.url),
			"secret":   t.secret,
		},
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	bodyReader := bytes.NewBuffer(jsonData)
	url := "https://api.twitch.tv/helix/eventsub/subscriptions"

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Client-Id", t.clientID)

	// Send the request
	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //nolint: errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//log.Println("subscribe response:", string(data))

	var s map[string]any
	err = json.Unmarshal(data, &s)
	if err != nil {
		return "", err
	}

	if _, ok := s["data"]; !ok {
		if msg, ok := s["message"]; ok {
			// to handle existing subscription
			return "", fmt.Errorf("%s", msg)
		}
		return "", fmt.Errorf("invalid data returned")
	}

	entries := s["data"].([]any)
	if len(entries) == 0 {
		return "", fmt.Errorf("no subscriptions returned")
	}

	sub := entries[0].(map[string]any)
	if _, ok := sub["id"]; !ok {
		return "", fmt.Errorf("no subscription id returned")
	}
	return sub["id"].(string), nil
}

// Subscriptions for a particular app
// requires an app access token
func (t *DefaultTwitchClient) subscriptions(token string) ([]*SubscriptionData, error) {
	url := "https://api.twitch.tv/helix/eventsub/subscriptions"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Client-Id", t.clientID)

	// Send the request
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() //nolint: errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var s SubscriptionResponse
	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	return s.Data, nil
}

// GetSubscription check if we have a subscription to the user
func (t *DefaultTwitchClient) GetSubscription(user string) (string, error) {
	token, err := t.GetAppAccessToken()
	if err != nil {
		return "", err
	}

	subs, err := t.subscriptions(token)
	if err != nil {
		return "", err
	}

	for _, sub := range subs {
		if u, ok := sub.Condition["broadcaster_user_id"]; ok && u == user {
			return sub.ID, nil
		}
	}
	return "", fmt.Errorf("subscription not found")
}

// SubscriptionResponse models data from twitch
type SubscriptionResponse struct {
	Total        int                 `json:"total"`
	Data         []*SubscriptionData `json:"data"`
	MaxTotalCost int                 `json:"max_total_cost"`
	TotalCost    int                 `json:"total_cost"`
	Pagination   any                 `json:"pagination"`
}

// SubscriptionData models data from twitch
type SubscriptionData struct {
	ID        string            `json:"id"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition map[string]string `json:"condition"`
	CreatedAt string            `json:"created_at"`
	Transport map[string]string `json:"transport"`
	Cost      int               `json:"cost"`
}
