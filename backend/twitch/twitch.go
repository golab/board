package twitch

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jarednogo/board/backend/core"
)

type Chat struct {
	Command string
	Body    string
}

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

// the base response type for twitch messages
type TwitchJSON struct {
	Subscription *Subscription    `json:"subscription"`
	Event        *TwitchEventJSON `json:"event"`
	Challenge    string           `json:"challenge"`
}

type Subscription struct {
	ID string `json:"id"`
}

// the twitch event type
type TwitchEventJSON struct {
	BroadcasterUserID string             `json:"broadcaster_user_id"`
	ChatterUserID     string             `json:"chatter_user_id"`
	Message           *TwitchMessageJSON `json:"message"`
}

// the twitch message type
type TwitchMessageJSON struct {
	Text string `json:"text"`
}

func Verify(message, signature string) bool {
	// get secret
	secret := Secret()
	if len(secret) == 0 {
		return true
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))

	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

func Secret() []byte {
	s := os.Getenv("TWITCHSECRET")
	return []byte(s)
}

func ClientID() string {
	return os.Getenv("TWITCHCLIENTID")
}

func BotID() string {
	return os.Getenv("TWITCHBOTID")
}

func GetUserAccessToken(code string) (string, error) {
	body := map[string]string{
		"client_id":     ClientID(),
		"client_secret": string(Secret()),
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  fmt.Sprintf("%s/apps/twitch/callback", core.MyURL()),
		"scope":         "channel:bot",
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	bodyReader := bytes.NewBuffer(jsonData)
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token")

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

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

func GetUsers(token string) (string, error) {
	url := "https://api.twitch.tv/helix/users"

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Client-Id", ClientID())

	client := &http.Client{Timeout: 10 * time.Second}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

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

func GetAppAccessToken() (string, error) {
	body := map[string]string{
		"client_id":     ClientID(),
		"client_secret": string(Secret()),
		"grant_type":    "client_credentials",
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	bodyReader := bytes.NewBuffer(jsonData)
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token")

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

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

type SubscriptionRequest struct {
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition map[string]string `json:"condition"`
	Transport map[string]string `json:"transport"`
}

func Subscribe(user, token string) (string, error) {
	body := SubscriptionRequest{
		Type:    "channel.chat.message",
		Version: "1",
		Condition: map[string]string{
			"broadcaster_user_id": user,
			"user_id":             BotID(),
		},
		Transport: map[string]string{
			"method":   "webhook",
			"callback": fmt.Sprintf("%s/apps/twitch/callback", core.MyURL()),
			"secret":   string(Secret()),
		},
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	bodyReader := bytes.NewBuffer(jsonData)
	url := fmt.Sprintf("https://api.twitch.tv/helix/eventsub/subscriptions")

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Client-Id", ClientID())

	client := &http.Client{Timeout: 10 * time.Second}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Println(string(data))

	var s map[string]interface{}
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

	entries := s["data"].([]interface{})
	if len(entries) == 0 {
		return "", fmt.Errorf("no subscriptions returned")
	}

	sub := entries[0].(map[string]interface{})
	if _, ok := sub["id"]; !ok {
		return "", fmt.Errorf("no subscription id returned")
	}
	return sub["id"].(string), nil
}
