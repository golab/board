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
	"log"
	// blank import is utilized to register the driver with the
	// database/sql package without directly using any of its
	// exported functions or types in the importing file
	_ "modernc.org/sqlite"
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
func Setup() {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return
	}
	defer db.Close()

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
		log.Println(err)
		return
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
		log.Println(err)
		return
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

	if err != nil {
		log.Println(err)
		return
	}
}

func TwitchGetRoom(broadcaster string) string {
	rooms, err := TwitchSelectRoom(broadcaster)
	if err != nil || len(rooms) != 1 {
		return ""
	}
	return rooms[0]
}

func TwitchSelectRoom(broadcaster string) ([]string, error) {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return []string{}, err
	}
	defer db.Close()
	rows, err := db.QueryContext(
		context.Background(),
		`SELECT roomid FROM twitch WHERE broadcaster = ?`,
		broadcaster)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	defer rows.Close()
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

func TwitchSetRoom(broadcaster, roomid string) error {
	rooms, err := TwitchSelectRoom(broadcaster)
	if err != nil {
		return err
	}

	if len(rooms) == 0 {
		return TwitchInsertRoom(broadcaster, roomid)
	} else {
		return TwitchUpdateRoom(broadcaster, roomid)
	}
}

func TwitchInsertRoom(broadcaster, roomid string) error {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.ExecContext(
		context.Background(),
		`INSERT INTO twitch (broadcaster, roomid) VALUES (?, ?)`,
		broadcaster, roomid)
	return err
}

func TwitchUpdateRoom(broadcaster, roomid string) error {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(
		context.Background(),
		`UPDATE twitch SET roomid = ? WHERE broadcaster = ?`,
		roomid, broadcaster)
	return err
}

func LoadRoom(id string) (*LoadJSON, error) {
	rows := SelectRoom(id)
	if len(rows) != 1 {
		return nil, fmt.Errorf("incorrect number of rows returned: %d", len(rows))
	}
	//fmt.Println("found:", id, rows[0])
	return rows[0], nil
}

// Save could also reasonably be called InsertOrUpdate
func SaveRoom(id string, data *LoadJSON) error {
	if len(SelectRoom(id)) != 0 {
		return UpdateRoom(id, data)
	} else {
		return InsertRoom(id, data)
	}
}

func DeleteRoom(id string) error {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(`DELETE FROM messages WHERE id = ?`, id)
	return err
}

func UpdateRoom(id string, data *LoadJSON) error {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return err
	}
	defer db.Close()
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

func InsertRoom(id string, data *LoadJSON) error {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return err
	}
	defer db.Close()

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

func LoadAllRooms() []*LoadJSON {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return []*LoadJSON{}
	}
	defer db.Close()

	rows, err := db.QueryContext(
		context.Background(),
		`SELECT id, sgf, loc, prefs, buffer, nextindex, password FROM rooms`)
	if err != nil {
		log.Println(err)
		return []*LoadJSON{}
	}

	defer rows.Close()
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

	return rooms
}

func SelectRoom(id string) []*LoadJSON {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return []*LoadJSON{}
	}
	defer db.Close()

	rows, err := db.QueryContext(
		context.Background(),
		`SELECT id, sgf, loc, prefs, buffer, nextindex, password FROM rooms WHERE id = ?`,
		id)
	if err != nil {
		log.Println(err)
		return []*LoadJSON{}
	}

	defer rows.Close()
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

	return rooms
}

func LoadAllMessages() []*MessageJSON {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return []*MessageJSON{}
	}
	defer db.Close()

	rows, err := db.QueryContext(
		context.Background(),
		`SELECT text, ttl FROM messages`)
	if err != nil {
		log.Println(err)
		return []*MessageJSON{}
	}

	defer rows.Close()
	var messages []*MessageJSON
	for rows.Next() {
		var text string
		var ttl int
		err = rows.Scan(&text, &ttl)
		msg := &MessageJSON{
			Text: text,
			TTL:  ttl,
		}
		messages = append(messages, msg)
	}

	return messages
}

func DeleteAllMessages() {
	db, err := sql.Open("sqlite", "file:"+Path())
	if err != nil {
		return
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM messages`)
}
