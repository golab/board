/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"path/filepath"
)

func Hash(input string) string {
	hashedBytes, _ := bcrypt.GenerateFromPassword(
		[]byte(input),
		bcrypt.DefaultCost)
	return string(hashedBytes)
}

func CorrectPassword(input, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	return err == nil
}

func CreateDir(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return false
		}
	}
	return true
}

func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	roomDir := filepath.Join(home, ".config", "tripleko")
	return roomDir
}

func Subpath(p string) string {
	return filepath.Join(Path(), p)
}

func RoomPath() string {
	return Subpath("rooms")
}

func MessagePath() string {
	return Subpath("messages")
}

func Setup() {
	roomDir := RoomPath()
	ok := CreateDir(roomDir)
	if !ok {
		log.Fatal("error creating room")
	}
	messageDir := MessagePath()
	ok = CreateDir(messageDir)
	if !ok {
		log.Fatal("error creating room")
	}
}

type LoadJSON struct {
	SGF       string         `json:"sgf"`
	Loc       string         `json:"loc"`
	Prefs     map[string]int `json:"prefs"`
	Buffer    int64          `json:"buffer"`
	NextIndex int            `json:"next_index"`
	Password  string         `json:"password"`
}

type EventJSON struct {
	Event  string      `json:"event"`
	Value  interface{} `json:"value"`
	Color  int         `json:"color"`
	UserID string      `json:"userid"`
}

func ErrorJSON(msg string) *EventJSON {
	return &EventJSON{"error", msg, 0, ""}
}

func FrameJSON(frame *Frame) *EventJSON {
	return &EventJSON{"frame", frame, 0, ""}
}

func NopJSON() *EventJSON {
	return &EventJSON{"nop", nil, 0, ""}
}
