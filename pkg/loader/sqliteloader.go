/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package loader

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	// blank import is utilized to register the driver with the
	// database/sql package without directly using any of its
	// exported functions or types in the importing file
	_ "modernc.org/sqlite"
)

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

type SqliteLoader struct {
	path string
	mu   sync.Mutex
}

func NewSqliteLoader(path string) *SqliteLoader {
	return &SqliteLoader{path: path}
}

func (ldr *SqliteLoader) Setup() error {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()
	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return err
	}
	defer db.Close() //nolint: errcheck

	// Create the 'rooms' table if it doesn't exist
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS rooms (
			id STRING PRIMARY KEY,
			sgf STRING NOT NULL,
			loc STRING NOT NULL,
			prefs STRING NOT NULL,
			buffer INTEGER NOT NULL,
			nextindex INTEGER NOT NULL,
			password STRING NOT NULL
		)`,
	)

	if err != nil {
		return err
	}

	// Create the 'messages' table if it doesn't exist
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text STRING NOT NULL,
			ttl INTEGER NOT NULL
		)`,
	)

	if err != nil {
		return err
	}

	// create the 'twitch' table if it doesn't exist
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS twitch (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			broadcaster STRING NOT NULL,
			roomid STRING
		)`,
	)

	return err
}

// Twitch logic

func (ldr *SqliteLoader) TwitchSetRoom(broadcaster, roomid string) error {
	rooms, err := ldr.twitchSelectRoom(broadcaster)
	if err != nil {
		return err
	}

	if len(rooms) == 0 {
		return ldr.twitchInsertRoom(broadcaster, roomid)
	} else {
		return ldr.twitchUpdateRoom(broadcaster, roomid)
	}
}

func (ldr *SqliteLoader) TwitchGetRoom(broadcaster string) string {
	rooms, err := ldr.twitchSelectRoom(broadcaster)
	if err != nil || len(rooms) != 1 {
		return ""
	}
	return rooms[0]
}

func (ldr *SqliteLoader) twitchSelectRoom(broadcaster string) ([]string, error) {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()

	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return []string{}, err
	}
	defer db.Close() //nolint: errcheck
	rows, err := db.QueryContext(
		context.Background(),
		`SELECT roomid FROM twitch WHERE broadcaster = ?`,
		broadcaster)
	if err != nil {
		return []string{}, err
	}

	defer rows.Close() //nolint: errcheck
	rooms := []string{}
	for rows.Next() {
		var roomid string
		err = rows.Scan(&roomid)
		if err != nil {
			continue
		}
		rooms = append(rooms, roomid)
	}

	return rooms, nil
}

func (ldr *SqliteLoader) twitchInsertRoom(broadcaster, roomid string) error {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()

	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return err
	}
	defer db.Close() //nolint: errcheck
	_, err = db.ExecContext(
		context.Background(),
		`INSERT INTO twitch (broadcaster, roomid) VALUES (?, ?)`,
		broadcaster, roomid)
	return err
}

func (ldr *SqliteLoader) twitchUpdateRoom(broadcaster, roomid string) error {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()

	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return err
	}
	defer db.Close() //nolint: errcheck

	_, err = db.ExecContext(
		context.Background(),
		`UPDATE twitch SET roomid = ? WHERE broadcaster = ?`,
		roomid, broadcaster)
	return err
}

func (ldr *SqliteLoader) LoadRoom(id string) (*LoadJSON, error) {
	rows, err := ldr.selectRoom(id)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, fmt.Errorf("incorrect number of rows returned: %d", len(rows))
	}
	return rows[0], nil
}

// Save could also reasonably be called InsertOrUpdate
func (ldr *SqliteLoader) SaveRoom(id string, data *LoadJSON) error {
	rows, err := ldr.selectRoom(id)
	if err != nil {
		return err
	}
	if len(rows) != 0 {
		return ldr.updateRoom(id, data)
	} else {
		return ldr.insertRoom(id, data)
	}
}

func (ldr *SqliteLoader) DeleteRoom(id string) error {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()

	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return err
	}
	defer db.Close() //nolint: errcheck
	_, err = db.Exec(`DELETE FROM rooms WHERE id = ?`, id)
	return err
}

func (ldr *SqliteLoader) updateRoom(id string, data *LoadJSON) error {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()

	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return err
	}
	defer db.Close() //nolint: errcheck
	prefs, err := Prefs(data.Prefs).ToString()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		context.Background(),
		`UPDATE rooms SET sgf = ?, loc = ?, prefs = ?, buffer = ?, nextindex = ?, password = ?  WHERE id = ?`,
		data.SGF, data.Location, prefs, data.Buffer, data.NextIndex, data.Password, id)
	return err
}

func (ldr *SqliteLoader) insertRoom(id string, data *LoadJSON) error {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()

	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return err
	}
	defer db.Close() //nolint: errcheck

	prefs, err := Prefs(data.Prefs).ToString()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		context.Background(),
		`INSERT INTO rooms (id, sgf, loc, prefs, buffer, nextindex, password) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, data.SGF, data.Location, prefs, data.Buffer, data.NextIndex, data.Password)

	return err
}

func (ldr *SqliteLoader) LoadAllRooms() ([]*LoadJSON, error) {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()

	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return nil, err
	}
	defer db.Close() //nolint: errcheck

	rows, err := db.QueryContext(
		context.Background(),
		`SELECT id, sgf, loc, prefs, buffer, nextindex, password FROM rooms`)
	if err != nil {
		return nil, err
	}

	defer rows.Close() //nolint: errcheck
	var rooms []*LoadJSON
	for rows.Next() {
		var id string
		var sgf string
		var loc string
		var prefsString string
		var buffer int
		var nextIndex int
		var password string
		err = rows.Scan(&id, &sgf, &loc, &prefsString, &buffer, &nextIndex, &password)
		if err != nil {
			return nil, err
		}
		prefs, _ := PrefsFromString(prefsString)
		data := &LoadJSON{
			SGF:       sgf,
			Location:  loc,
			Prefs:     prefs,
			Buffer:    int64(buffer),
			NextIndex: nextIndex,
			Password:  password,
			ID:        id,
		}
		rooms = append(rooms, data)
	}

	return rooms, nil
}

func (ldr *SqliteLoader) selectRoom(id string) ([]*LoadJSON, error) {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()

	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return nil, err
	}
	defer db.Close() //nolint: errcheck

	rows, err := db.QueryContext(
		context.Background(),
		`SELECT id, sgf, loc, prefs, buffer, nextindex, password FROM rooms WHERE id = ?`,
		id)
	if err != nil {
		return nil, err
	}

	defer rows.Close() //nolint: errcheck
	var rooms []*LoadJSON
	for rows.Next() {
		var id string
		var sgf string
		var loc string
		var prefsString string
		var buffer int
		var nextIndex int
		var password string
		err = rows.Scan(&id, &sgf, &loc, &prefsString, &buffer, &nextIndex, &password)
		if err != nil {
			return nil, err
		}
		prefs, _ := PrefsFromString(prefsString)
		data := &LoadJSON{
			SGF:       sgf,
			Location:  loc,
			Prefs:     prefs,
			Buffer:    int64(buffer),
			NextIndex: nextIndex,
			Password:  password,
			ID:        id,
		}
		rooms = append(rooms, data)
	}

	return rooms, nil
}

func (ldr *SqliteLoader) LoadAllMessages() ([]*MessageJSON, error) {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()

	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return nil, err
	}
	defer db.Close() //nolint: errcheck

	rows, err := db.QueryContext(
		context.Background(),
		`SELECT text, ttl FROM messages`)
	if err != nil {
		return nil, err
	}

	defer rows.Close() //nolint: errcheck
	var messages []*MessageJSON
	for rows.Next() {
		var text string
		var ttl int
		err = rows.Scan(&text, &ttl)
		if err != nil {
			return nil, err
		}
		msg := &MessageJSON{
			Text: text,
			TTL:  ttl,
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (ldr *SqliteLoader) DeleteAllMessages() error {
	ldr.mu.Lock()
	defer ldr.mu.Unlock()

	db, err := sql.Open("sqlite", "file:"+ldr.path)
	if err != nil {
		return err
	}
	defer db.Close() //nolint: errcheck

	_, err = db.Exec(`DELETE FROM messages`)
	return err
}
