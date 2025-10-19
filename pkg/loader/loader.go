/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package loader

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type MessageJSON struct {
	Text string `json:"message"`
	TTL  int    `json:"ttl"`
}

type LoadJSON struct {
	SGF       string         `json:"sgf"`
	Location  string         `json:"loc"`
	Prefs     map[string]int `json:"prefs"`
	Buffer    int64          `json:"buffer"`
	NextIndex int            `json:"next_index"`
	Password  string         `json:"password"`
	ID        string         `json:"id"`
}

type Prefs map[string]int

func (p Prefs) ToString() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func PrefsFromString(s string) (Prefs, error) {
	p := make(map[string]int)
	err := json.Unmarshal([]byte(s), &p)
	if err != nil {
		return nil, err
	}
	return Prefs(p), nil
}

func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	newpath := filepath.Join(home, ".config", "tripleko")
	err = os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		home = "."
	}

	dbPath := filepath.Join(home, ".config", "tripleko", "board.db")
	return dbPath
}

type Loader interface {
	Setup()
	TwitchGetRoom(string) string
	TwitchSetRoom(string, string) error

	SaveRoom(string, *LoadJSON) error
	LoadRoom(string) (*LoadJSON, error)
	LoadAllRooms() []*LoadJSON
	DeleteRoom(string) error
	LoadAllMessages() []*MessageJSON
	DeleteAllMessages()
}

func NewDefaultLoader() Loader {
	return NewSqliteLoader()
}
