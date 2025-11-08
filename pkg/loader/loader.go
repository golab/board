/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package loader

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

type Loader interface {
	Setup() error
	TwitchGetRoom(string) string
	TwitchSetRoom(string, string) error

	SaveRoom(string, *LoadJSON) error
	LoadRoom(string) (*LoadJSON, error)
	LoadAllRooms() ([]*LoadJSON, error)
	DeleteRoom(string) error
	LoadAllMessages() ([]*MessageJSON, error)
	DeleteAllMessages() error
}

func NewDefaultLoader(path string) Loader {
	return NewSqliteLoader(path)
}
