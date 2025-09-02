package twitch

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

type Chat struct {
	Command string
	Body    string
}

func Parse(chat string) (*Chat, error) {
	chat = strings.TrimSpace(chat)
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
	Subscription interface{}      `json:"subscription"`
	Event        *TwitchEventJSON `json:"event"`
	Challenge    string           `json:"challenge"`
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
	secret := GetTwitchSecret()
	if len(secret) == 0 {
		return true
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))

	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

func GetTwitchSecret() []byte {
	s := os.Getenv("TWITCHSECRET")
	return []byte(s)
}
