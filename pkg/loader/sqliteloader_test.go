/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package loader_test

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/pkg/loader"
)

func TestSqliteLoader(t *testing.T) {
	// make a temporary dir to store the db
	tmp := t.TempDir()
	path := filepath.Join(tmp, "board.db")

	// test setup
	ldr, err := loader.NewSqliteLoader(path)
	assert.NoError(t, err)

	// test adding a twitch room
	err = ldr.TwitchSetRoom("user", "board")
	assert.NoError(t, err)

	// test loading the same room
	b := ldr.TwitchGetRoom("user")
	assert.Equal(t, b, "board")

	// test updating the same room
	err = ldr.TwitchSetRoom("user", "otherboard")
	assert.NoError(t, err)

	// test loading the same room
	b = ldr.TwitchGetRoom("user")
	assert.Equal(t, b, "otherboard")

	// make a struct to save a room
	l := &loader.LoadJSON{
		SGF:       "test",
		Location:  "test",
		Buffer:    500,
		NextIndex: 42,
		Password:  "test",
		ID:        "test",
	}

	// test saving a room
	err = ldr.SaveRoom(l.ID, l)
	assert.NoError(t, err)

	// test loading the same room
	m, err := ldr.LoadRoom(l.ID)
	assert.NoError(t, err)
	assert.Equal(t, m.SGF, l.SGF)
	assert.Equal(t, m.Location, l.Location)
	assert.Equal(t, m.Buffer, l.Buffer)
	assert.Equal(t, m.NextIndex, l.NextIndex)
	assert.Equal(t, m.Password, l.Password)
	assert.Equal(t, m.ID, l.ID)

	// test saving a new room
	m.ID = "test2"
	err = ldr.SaveRoom(m.ID, m)
	assert.NoError(t, err)

	// test loading all the rooms
	rooms, err := ldr.LoadAllRooms()
	assert.NoError(t, err)
	assert.Equal(t, len(rooms), 2)

	// test updating a room
	err = ldr.SaveRoom(m.ID, m)
	assert.NoError(t, err)

	// test deleting a room
	err = ldr.DeleteRoom(m.ID)
	assert.NoError(t, err)
	_, err = ldr.LoadRoom(m.ID)
	assert.NotNil(t, err)

	// test loading messages
	msgs, err := ldr.LoadAllMessages()
	assert.NoError(t, err)
	assert.Equal(t, len(msgs), 0)

	// InsertTestMessage is only for testing
	err = ldr.InsertTestMessage()
	assert.NoError(t, err)

	// test loading all messages
	msgs, err = ldr.LoadAllMessages()
	assert.NoError(t, err)
	assert.Equal(t, len(msgs), 1)

	// test deleting all messages
	err = ldr.DeleteAllMessages()
	assert.NoError(t, err)

	// test loading messages again
	msgs, err = ldr.LoadAllMessages()
	assert.NoError(t, err)
	assert.Equal(t, len(msgs), 0)
}

func TestSqliteLoaderConcurrency(t *testing.T) {
	// make a temporary dir to store the db
	tmp := t.TempDir()
	path := filepath.Join(tmp, "board.db")

	// test setup
	ldr, err := loader.NewSqliteLoader(path)
	assert.NoError(t, err)

	// create wait group
	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			// make a struct to save a room
			l := &loader.LoadJSON{
				SGF:       "test",
				Location:  "test",
				Buffer:    500,
				NextIndex: 42,
				Password:  "test",
				ID:        fmt.Sprintf("%d", i),
			}

			// test saving a room
			err := ldr.SaveRoom(l.ID, l)
			assert.NoError(t, err)
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestSqliteLoaderAutoNumber(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "board.db")

	ldr, err := loader.NewSqliteLoader(path)
	assert.NoError(t, err)

	l := &loader.LoadJSON{SGF: "test", ID: "test", AutoNumber: true}
	assert.NoError(t, ldr.SaveRoom(l.ID, l))
	m, err := ldr.LoadRoom(l.ID)
	assert.NoError(t, err)
	assert.True(t, m.AutoNumber)

	// update path
	m.AutoNumber = false
	assert.NoError(t, ldr.SaveRoom(m.ID, m))
	m, err = ldr.LoadRoom(l.ID)
	assert.NoError(t, err)
	assert.False(t, m.AutoNumber)
}

func TestSqliteLoaderMigrateAutoNumber(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "board.db")

	// a database created before the autonumber column existed
	db, err := sql.Open("sqlite", "file:"+path)
	assert.NoError(t, err)
	_, err = db.Exec(
		`CREATE TABLE rooms (
			id TEXT PRIMARY KEY,
			sgf TEXT NOT NULL,
			loc TEXT NOT NULL,
			prefs TEXT NOT NULL,
			buffer INTEGER NOT NULL,
			nextindex INTEGER NOT NULL,
			password TEXT NOT NULL
		)`)
	assert.NoError(t, err)
	_, err = db.Exec(
		`INSERT INTO rooms (id, sgf, loc, prefs, buffer, nextindex, password) VALUES ('old', 'sgf', '0', '{}', 250, 1, '')`)
	assert.NoError(t, err)
	assert.NoError(t, db.Close())

	ldr, err := loader.NewSqliteLoader(path)
	assert.NoError(t, err)

	m, err := ldr.LoadRoom("old")
	assert.NoError(t, err)
	assert.Equal(t, m.SGF, "sgf")
	assert.False(t, m.AutoNumber)

	m.AutoNumber = true
	assert.NoError(t, ldr.SaveRoom(m.ID, m))
	m, err = ldr.LoadRoom("old")
	assert.NoError(t, err)
	assert.True(t, m.AutoNumber)

	// running setup again on an already-migrated db is fine
	assert.NoError(t, ldr.Close())
	_, err = loader.NewSqliteLoader(path)
	assert.NoError(t, err)
}
