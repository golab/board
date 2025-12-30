/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package hub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golab/board/internal/twitch"
	"github.com/golab/board/pkg/core"
	"github.com/golab/board/pkg/event"
	"github.com/google/uuid"
)

func (h *Hub) TwitchRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/subscribe", h.twitchSubscribe)
	r.Get("/unsubscribe", h.twitchUnsubscribe)
	r.Get("/callback", h.twitchCallbackGet)
	r.Post("/callback", h.twitchCallbackPost)
	return r
}

func (h *Hub) twitchSubscribe(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	expiration := time.Now().Add(2 * time.Minute)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   true,
		Expires:  expiration,
		Path:     "/",
	})
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s/apps/twitch/callback&scope=%s&state=%s", h.cfg.Twitch.ClientID, h.cfg.Server.URL, "channel:bot", state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Hub) twitchUnsubscribe(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	expiration := time.Now().Add(2 * time.Minute)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   true,
		Expires:  expiration,
		Path:     "/",
	})
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s/apps/twitch/callback&state=%s", h.cfg.Twitch.ClientID, h.cfg.Server.URL, state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Hub) twitchCallbackGet(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")

	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != state {
		http.Error(w, "invalid state", http.StatusForbidden)
		return
	}

	if code != "" {

		// use the code to get an access token
		token, err := h.tc.GetUserAccessToken(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		// use the user access token to get the user id
		user, err := h.tc.GetUsers(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		// get an app access token (one could imagine putting this
		// in the subscribe function directly)
		token, err = h.tc.GetAppAccessToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if scope == "" {
			// unsubscribe logic
			id, err := h.tc.GetSubscription(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			// unsubscribe
			err = h.tc.Unsubscribe(id, token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			h.logger.Info("unsubscribing", "id", id, "user", user)
		} else {
			// subscribe, get subscription id
			id, err := h.tc.Subscribe(user, token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			h.logger.Info("new subscription", "id", id, "user", user)
		}
	}

	//w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte(`{"message": "success"}`))
	_, err = w.Write([]byte("success"))
	if err != nil {
		h.logger.Error("twitchCallbackGet:write", "err", err)
	}
}

func (h *Hub) twitchCallbackPost(w http.ResponseWriter, r *http.Request) {
	// read the body into a []byte
	body, _ := io.ReadAll(r.Body)

	// try to read the body into a TwitchJSON struct
	var req twitch.TwitchJSON
	err := json.Unmarshal(body, &req)
	if err != nil {
		h.logger.Error("twitchCallbackPost:unmarshal", "err", err)
		return
	}

	// on subscriptions, twitch sends a challenge that we need to respond to
	if req.Challenge != "" {
		_, err = w.Write([]byte(req.Challenge))
		if err != nil {
			h.logger.Error("twitchCallbackPost:write", "err", err)
		}
		return
	}

	// Grab headers for verification
	msgid := r.Header.Get("Twitch-Eventsub-Message-Id")
	timestamp := r.Header.Get("Twitch-Eventsub-Message-Timestamp")
	signature := r.Header.Get("Twitch-Eventsub-Message-Signature")

	// concat for verification
	message := msgid + timestamp + string(body)

	// do verification
	v := twitch.Verify(h.cfg.Twitch.Secret, message, signature)
	if !v {
		h.logger.Info("unverified message")
		return
	}

	// try to pull the subscription
	subsc := req.Subscription
	if subsc != nil {
		h.logger.Debug("subscription parsed", "subscription", subsc)
	}

	// try to pull out the event
	evt := req.Event
	if evt == nil {
		h.logger.Info("no event parsed")
		return
	}

	// try to pull out the message
	if evt.Message == nil {
		h.logger.Info("no message parsed")
		return
	}

	// get broadcaster and chatter
	broadcaster := evt.BroadcasterUserID
	chatter := evt.ChatterUserID

	// extract the message in chat
	text := evt.Message.Text

	if !strings.HasPrefix(text, "!") {
		h.logger.Debug("twitch message", "text", text, "broadcaster", broadcaster)
		return
	}

	chat, err := twitch.Parse(text)
	if err != nil {
		h.logger.Error("twitchCallBackPost:parse", "err", err, "broadcaster", broadcaster)
		return
	}

	// only care about the relevant commands
	if chat.Command != "branch" && chat.Command != "setboard" {
		h.logger.Info("invalid command", "command", chat.Command, "broadcaster", broadcaster)
		return
	}

	h.logger.Info("received", "command", chat.Command, "body", chat.Body, "broadcaster", broadcaster)

	// make sure only the broadcaster can set the room
	switch chat.Command {
	case "setboard":
		if broadcaster == chatter {
			tokens := strings.Split(chat.Body, " ")
			if len(tokens) == 0 {
				return
			}
			roomID := core.Sanitize(tokens[0])
			if len(roomID) == 0 {
				h.logger.Info("empty roomID", "broadcaster", broadcaster)
				return
			}

			h.logger.Info("setting roomid", "broadcaster", broadcaster, "room_id", roomID)
			err := h.db.TwitchSetRoom(broadcaster, roomID)
			if err != nil {
				h.logger.Error("error in setboard", "err", err, "broadcaster", broadcaster)
			}
		} else {
			h.logger.Info("unauthorized user tried to setboard", "broadcaster", broadcaster)
		}
	case "branch":
		branch := strings.ToLower(chat.Body)

		roomID := h.db.TwitchGetRoom(broadcaster)
		if roomID == "" {
			h.logger.Info("room not set", "broadcaster", broadcaster)
			return
		}
		h.logger.Info("room found", "broadcaster", broadcaster, "room_id", roomID)
		r := h.GetOrCreateRoom(roomID)

		// create the event
		e := event.NewEvent("graft", branch)

		h.logger.Info("grafting", "branch", branch, "broadcaster", broadcaster, "room_id", roomID)
		r.HandleAny(e)
	}

	w.WriteHeader(http.StatusOK)
}
