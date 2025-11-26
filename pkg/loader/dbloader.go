/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package loader

import (
	"database/sql"
	"fmt"
	"strings"
)

type DBType int

const (
	DBTypeSqlite DBType = iota
	DBTypePostgres
)

type BaseLoader struct {
	db     *sql.DB
	dbType DBType
}

func (ldr *BaseLoader) Exec(query string, args ...interface{}) (sql.Result, error) {
	q := ldr.rebind(query)
	return ldr.db.Exec(q, args...)
}

func (ldr *BaseLoader) Query(query string, args ...interface{}) (*sql.Rows, error) {
	q := ldr.rebind(query)
	return ldr.db.Query(q, args...)
}

func (ldr *BaseLoader) setup() error {
	switch ldr.dbType {
	case DBTypeSqlite:
		return ldr.sqliteSetup()
	case DBTypePostgres:
		return ldr.postgresSetup()
	}
	return fmt.Errorf("invalid dbtype")
}

func (ldr *BaseLoader) sqliteSetup() error {
	// Create the 'rooms' table if it doesn't exist
	_, err := ldr.db.Exec(
		`CREATE TABLE IF NOT EXISTS rooms (
			id TEXT PRIMARY KEY,
			sgf TEXT NOT NULL,
			loc TEXT NOT NULL,
			prefs TEXT NOT NULL,
			buffer INTEGER NOT NULL,
			nextindex INTEGER NOT NULL,
			password TEXT NOT NULL
		)`,
	)

	if err != nil {
		return err
	}

	// Create the 'messages' table if it doesn't exist
	_, err = ldr.db.Exec(
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT NOT NULL,
			ttl INTEGER NOT NULL
		)`,
	)

	if err != nil {
		return err
	}

	// create the 'twitch' table if it doesn't exist
	_, err = ldr.db.Exec(
		`CREATE TABLE IF NOT EXISTS twitch (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			broadcaster TEXT NOT NULL,
			roomid TEXT
		)`,
	)

	return err
}

func (ldr *BaseLoader) postgresSetup() error {
	// Create the 'rooms' table if it doesn't exist
	_, err := ldr.db.Exec(
		`CREATE TABLE IF NOT EXISTS rooms (
			id TEXT PRIMARY KEY,
			sgf TEXT NOT NULL,
			loc TEXT NOT NULL,
			prefs TEXT NOT NULL,
			buffer INTEGER NOT NULL,
			nextindex INTEGER NOT NULL,
			password TEXT NOT NULL
		)`,
	)

	if err != nil {
		return err
	}

	// Create the 'messages' table if it doesn't exist
	_, err = ldr.db.Exec(
		`CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			text TEXT NOT NULL,
			ttl INTEGER NOT NULL
		)`,
	)

	if err != nil {
		return err
	}

	// create the 'twitch' table if it doesn't exist
	_, err = ldr.db.Exec(
		`CREATE TABLE IF NOT EXISTS twitch (
			id SERIAL PRIMARY KEY,
			broadcaster TEXT NOT NULL,
			roomid TEXT
		)`,
	)

	return err
}

// rebind converts '?' -> '$1' '$2' ... when using DBTypePostgres
func (ldr *BaseLoader) rebind(query string) string {
	if ldr.dbType == DBTypeSqlite {
		return query
	}

	// DollarStyle: replace each '?' with $1, $2, ...
	var b strings.Builder
	argIdx := 1
	for i := 0; i < len(query); i++ {
		ch := query[i]
		if ch == '?' {
			b.WriteString(fmt.Sprintf("$%d", argIdx))
			argIdx++
		} else {
			b.WriteByte(ch)
		}
	}
	return b.String()
}

func (ldr *BaseLoader) Close() error {
	return ldr.db.Close()
}

// Twitch logic

func (ldr *BaseLoader) TwitchSetRoom(broadcaster, roomid string) error {
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

func (ldr *BaseLoader) TwitchGetRoom(broadcaster string) string {
	rooms, err := ldr.twitchSelectRoom(broadcaster)
	if err != nil || len(rooms) != 1 {
		return ""
	}
	return rooms[0]
}

func (ldr *BaseLoader) twitchSelectRoom(broadcaster string) ([]string, error) {
	rows, err := ldr.Query(
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

func (ldr *BaseLoader) twitchInsertRoom(broadcaster, roomid string) error {
	_, err := ldr.Exec(
		`INSERT INTO twitch (broadcaster, roomid) VALUES (?, ?)`,
		broadcaster, roomid)
	return err
}

func (ldr *BaseLoader) twitchUpdateRoom(broadcaster, roomid string) error {
	_, err := ldr.Exec(
		`UPDATE twitch SET roomid = ? WHERE broadcaster = ?`,
		roomid, broadcaster)
	return err
}

func (ldr *BaseLoader) LoadRoom(id string) (*LoadJSON, error) {
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
func (ldr *BaseLoader) SaveRoom(id string, data *LoadJSON) error {
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

func (ldr *BaseLoader) DeleteRoom(id string) error {
	_, err := ldr.Exec(`DELETE FROM rooms WHERE id = ?`, id)
	return err
}

func (ldr *BaseLoader) updateRoom(id string, data *LoadJSON) error {
	prefs, err := Prefs(data.Prefs).ToString()
	if err != nil {
		return err
	}

	_, err = ldr.Exec(
		`UPDATE rooms SET sgf = ?, loc = ?, prefs = ?, buffer = ?, nextindex = ?, password = ?  WHERE id = ?`,
		data.SGF, data.Location, prefs, data.Buffer, data.NextIndex, data.Password, id)
	return err
}

func (ldr *BaseLoader) insertRoom(id string, data *LoadJSON) error {
	prefs, err := Prefs(data.Prefs).ToString()
	if err != nil {
		return err
	}

	_, err = ldr.Exec(
		`INSERT INTO rooms (id, sgf, loc, prefs, buffer, nextindex, password) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, data.SGF, data.Location, prefs, data.Buffer, data.NextIndex, data.Password)

	return err
}

func (ldr *BaseLoader) LoadAllRooms() ([]*LoadJSON, error) {
	rows, err := ldr.Query(
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

func (ldr *BaseLoader) selectRoom(id string) ([]*LoadJSON, error) {
	rows, err := ldr.Query(
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

func (ldr *BaseLoader) LoadAllMessages() ([]*MessageJSON, error) {
	rows, err := ldr.Query(
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

func (ldr *BaseLoader) DeleteAllMessages() error {
	_, err := ldr.Exec(`DELETE FROM messages`)
	return err
}

func (ldr *BaseLoader) InsertTestMessage() error {
	_, err := ldr.Exec(`INSERT INTO messages (text, ttl) VALUES ("foo", 60)`)
	return err
}
