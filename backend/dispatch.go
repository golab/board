/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"encoding/base64"
	"log"
	"strconv"
	"strings"
	"time"
)

type EventHandler func(*EventJSON) *EventJSON

type Middleware func(EventHandler) EventHandler


// handlers

func (room *Room) HandleIsProtected(evt *EventJSON) *EventJSON {
	evt.Value = room.HasPassword()
	room.SendTo(evt.UserID, evt)
	return evt
}

func (room *Room) HandleCheckPassword(evt *EventJSON) *EventJSON {
	p := evt.Value.(string)

	if !CorrectPassword(p, room.password) {
		evt.Value = ""
	} else {
		room.auth[evt.UserID] = true
	}
	room.SendTo(evt.UserID, evt)
	return evt
}

func HandleDebug(evt *EventJSON) *EventJSON {
	log.Println(evt.UserID, evt)
	return evt
}

func HandlePing(evt *EventJSON) *EventJSON {
	return evt
}

func (room *Room) HandleUploadSGF(evt *EventJSON) *EventJSON {
	var bcast *EventJSON

	// it might be a string
	if str, ok := evt.Value.(string); ok {

		decoded, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			bcast = ErrorJSON(err.Error())
		}
		if IsZipFile(decoded) {
			filesBytes, err := Decompress(decoded)
			if err != nil {
				bcast = ErrorJSON(err.Error())
			} else {
				files := []string{}
				for _, file := range filesBytes {
					files = append(files, string(file))
				}
				merged := Merge(files)
				bcast = room.UploadSGF(merged)
			}
		} else {
			bcast = room.UploadSGF(string(decoded))
		}

	} else if arr, ok := evt.Value.([]interface{}); ok {
		// it might be an array of strings
		sgfs := []string{}
		for _, ifc := range arr {
			str := ifc.(string)
			d, err := base64.StdEncoding.DecodeString(str)
			if err != nil {
				bcast = ErrorJSON(err.Error())
			}
			sgfs = append(sgfs, string(d))
		}
		sgf := Merge(sgfs)
		bcast = room.UploadSGF(sgf)
	} else {
		bcast = ErrorJSON("unreachable")
	}

	bcast.UserID = evt.UserID
	return bcast
}

func (room *Room) HandleRequestSGF(evt *EventJSON) *EventJSON {
	var bcast *EventJSON
	defer func(){if bcast != nil {bcast.UserID = evt.UserID}}()

	url := evt.Value.(string)
	
	if IsOGS(url) {

		connectToOGS := false

		spl := strings.Split(url, "/")
		if len(spl) < 2 {
			bcast = ErrorJSON("url parsing error")
			return bcast
		}

		ogsType := spl[len(spl)-2]
		//log.Println(ogsType)

		if ogsType == "game" {
			ended, err := OGSCheckEnded(url)
			if err != nil {
				bcast = ErrorJSON(err.Error())
				return bcast
			}
			connectToOGS = !ended
		} else if ogsType == "review" || ogsType == "demo"   {
			ogsType = "review" 
			connectToOGS = true
		}

		//log.Println(connectToOGS)
		if connectToOGS {

			idStr := spl[len(spl)-1]
			id64, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				bcast = ErrorJSON("int parsing error")
				return bcast
			}
			id := int(id64)

			o, err := NewOGSConnector(room)
			if err != nil {
				bcast = ErrorJSON("ogs connector error")
				return bcast
			}
			go o.Loop(id,ogsType)
			room.OGSLink = o

			if(ogsType == "game"){
				// finish here
				return NopJSON()
			}
		}
	}

	data, err := ApprovedFetch(evt.Value.(string))
	if err != nil {
		bcast = ErrorJSON(err.Error())
	} else if data == "Permission denied" {
		bcast = ErrorJSON("Error fetching SGF. Is it a private OGS game?")
	} else {
		bcast = room.UploadSGF(string(data))
	}

	return bcast
}

func (room *Room) HandleTrash(evt *EventJSON) *EventJSON {

	// reset room
	oldBuffer := room.State.InputBuffer
	room.State = NewState(room.State.Size, true)

	// reuse old inputbuffer
	room.State.InputBuffer = oldBuffer

	frame := room.State.GenerateFullFrame(true)
	bcast := FrameJSON(frame)
	bcast.UserID = evt.UserID
	return bcast
}

func (room *Room) HandleUpdateNickname(evt *EventJSON) *EventJSON {
	nickname := evt.Value.(string)
	room.nicks[evt.UserID] = nickname
	userEvt := &EventJSON{
		"connected_users",
		room.nicks,
		0,
		evt.UserID,
	}
	return userEvt
}

func (room *Room) HandleUpdateSettings(evt *EventJSON) *EventJSON {
	sMap := evt.Value.(map[string]interface{})
	buffer := int64(sMap["buffer"].(float64))
	size := int(sMap["size"].(float64))
	nickname := sMap["nickname"].(string)

	room.nicks[evt.UserID] = nickname

	password := sMap["password"].(string)
	hashed := ""
	if password != "" {
		hashed = Hash(password)
	}
	settings := &Settings{buffer, size, hashed}

	room.State.InputBuffer = settings.Buffer
	if settings.Size != room.State.Size {
		// essentially trashing
		room.State = NewState(settings.Size, true)
		room.State.InputBuffer = buffer
	}

	// can be changed
	// anyone already in the room is added
	// person who set password automatically gets added
	for connID, _ := range room.conns {
		room.auth[connID] = true
	}
	room.password = hashed

	return evt
}

func (room *Room) HandleEvent(evt *EventJSON) *EventJSON {
	var bcast *EventJSON
	frame, err := room.State.AddEvent(evt)
	if err != nil {
		bcast = ErrorJSON(err.Error())
	}
	if frame != nil {
		bcast = FrameJSON(frame)
	} else {
		bcast = evt
	}
	bcast.UserID = evt.UserID
	return bcast
}

// middleware

func (room *Room) BroadcastAfter(setTime bool) Middleware {
	return func(handler EventHandler) EventHandler {
		return func(evt *EventJSON) *EventJSON {
			evt = handler(evt)
			room.Broadcast(evt, setTime)
			return evt
		}
	}
}

func (room *Room) BroadcastFullFrameAfter(handler EventHandler) EventHandler {
	return func(evt *EventJSON) *EventJSON {
		evt = handler(evt)
		frame := room.State.GenerateFullFrame(true)
		bcast := FrameJSON(frame)
		room.Broadcast(bcast, false)
		return evt
	}
}

func (room *Room) BroadcastConnectedUsersAfter(handler EventHandler) EventHandler {
	return func(evt *EventJSON) *EventJSON {
		evt = handler(evt)
		userEvt := &EventJSON{
			"connected_users",
			room.nicks,
			0,
			"",
		}

		// broadcast connected_users
		room.Broadcast(userEvt, false)
		return evt
	}
}

func (room *Room) CloseOGS(handler EventHandler) EventHandler {
	return func(evt *EventJSON) *EventJSON {
		if room.OGSLink != nil {
			room.OGSLink.End()
		}
		return handler(evt)
	}
}

func (room *Room) Authorized(handler EventHandler) EventHandler {
	return func(evt *EventJSON) *EventJSON {
		id := evt.UserID
		_, ok := room.auth[id]
		if room.password == "" || ok {
			// only go to the next handler if authorized
			evt = handler(evt)
		}
		return evt
	}
}

// this one is to keep the same user from submitting multiple events too quickly
func (room *Room) Slow(handler EventHandler) EventHandler {
	return func(evt *EventJSON) *EventJSON {
		id := evt.UserID
		// check multiple events from the same user in a narrow window (50 ms)
		now := time.Now()
		if last, ok := room.lastMessages[id]; !ok {
			room.lastMessages[id] = &now
		} else {
			diff := now.Sub(*last)
			room.lastMessages[id] = &now
			if diff.Milliseconds() < 50 {
				// don't do the next handler if too fast
				return evt
			}
		}
		return handler(evt)
	}
}

// this one is to keep people from tripping over each other
func (room *Room) OutsideBuffer(handler EventHandler) EventHandler {
	return func(evt *EventJSON) *EventJSON {
		if room.lastUser != evt.UserID {
			now := time.Now()
			diff := now.Sub(*room.timeLastEvent)
			if diff.Milliseconds() < room.State.InputBuffer {
				// don't do the next handler if too fast
				return evt
			}
		}
		return handler(evt)
	}
}

func Chain(h EventHandler, middleware ...Middleware) EventHandler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}
