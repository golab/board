/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package loader_test

import (
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/loader"
)

func TestMemoryLoader(t *testing.T) {
	ml := loader.NewMemoryLoader()

	ml.AddMessage("test", 30)
	assert.Equal(t, ml.MessageCount(), 1)

	msgs, err := ml.LoadAllMessages()
	assert.Equal(t, len(msgs), 1)
	assert.NoError(t, err)

	err = ml.DeleteAllMessages()
	assert.NoError(t, err)

	msgs, err = ml.LoadAllMessages()
	assert.Equal(t, len(msgs), 0)
	assert.NoError(t, err)

	assert.Equal(t, ml.TwitchGetRoom("user123"), "")
	err = ml.TwitchSetRoom("user123", "room1")
	assert.NoError(t, err)
	assert.Equal(t, ml.TwitchGetRoom("user123"), "room1")

	// make a struct to save a room
	l := &loader.LoadJSON{
		SGF:       "test",
		Location:  "test",
		Buffer:    500,
		NextIndex: 42,
		Password:  "test",
		ID:        "test",
	}

	_, err = ml.LoadRoom("room1")
	assert.NotNil(t, err)
	err = ml.SaveRoom("room1", l)
	assert.NoError(t, err)
	m, err := ml.LoadRoom("room1")
	assert.Equal(t, m.SGF, l.SGF)
	assert.NoError(t, err)

	rooms, err := ml.LoadAllRooms()
	assert.Equal(t, len(rooms), 1)
	assert.NoError(t, err)

	err = ml.DeleteRoom("room1")
	assert.NoError(t, err)
	rooms, err = ml.LoadAllRooms()
	assert.Equal(t, len(rooms), 0)
	assert.NoError(t, err)

}
