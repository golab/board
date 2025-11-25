/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package loader

import (
	"fmt"
)

type MemoryLoader struct {
	rooms    map[string]*LoadJSON
	messages []*MessageJSON
	twitch   map[string]string
}

func NewMemoryLoader() *MemoryLoader {
	return &MemoryLoader{
		rooms:  make(map[string]*LoadJSON),
		twitch: make(map[string]string),
	}
}

func (ml *MemoryLoader) AddMessage(text string, ttl int) {
	ml.messages = append(ml.messages, &MessageJSON{text, ttl})
}

func (ml *MemoryLoader) MessageCount() int {
	return len(ml.messages)
}

func (ml *MemoryLoader) TwitchGetRoom(broadcaster string) string {
	if roomID, ok := ml.twitch[broadcaster]; ok {
		return roomID
	}
	return ""
}

func (ml *MemoryLoader) TwitchSetRoom(broadcaster, roomID string) error {
	ml.twitch[broadcaster] = roomID

	return nil
}

func (ml *MemoryLoader) SaveRoom(s string, l *LoadJSON) error {
	ml.rooms[s] = l
	return nil
}

func (ml *MemoryLoader) LoadRoom(s string) (*LoadJSON, error) {
	if l, ok := ml.rooms[s]; ok {
		return l, nil
	}
	return nil, fmt.Errorf("room not found")
}

func (ml *MemoryLoader) LoadAllRooms() ([]*LoadJSON, error) {
	ls := []*LoadJSON{}
	for _, l := range ml.rooms {
		ls = append(ls, l)
	}
	return ls, nil
}

func (ml *MemoryLoader) DeleteRoom(s string) error {
	delete(ml.rooms, s)
	return nil
}

func (ml *MemoryLoader) LoadAllMessages() ([]*MessageJSON, error) {
	return ml.messages, nil
}

func (ml *MemoryLoader) DeleteAllMessages() error {
	ml.messages = []*MessageJSON{}
	return nil
}

func (ml *MemoryLoader) Close() error {
	return nil
}
