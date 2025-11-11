/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package plugin

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/event"
	"github.com/jarednogo/board/pkg/fetch"
	"golang.org/x/net/websocket"
)

type Room interface {
	HeadColor() core.Color
	PushHead(int, int, core.Color)
	GenerateFullFrame(core.TreeJSONType) *core.Frame
	AddStones([]*core.Stone)
	Broadcast(event.Event)
	UploadSGF(string) event.Event
}

/*
func GetUser(f fetch.Fetcher, id int) (string, error) {
	data, err := f.Fetch(fmt.Sprintf("http://online-go.com/api/v1/players/%d/", id))
	if err != nil {
		return "", err
	}
	log.Println(string(data))
	user := &User{}
	err = json.Unmarshal([]byte(data), user)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}
*/

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Creds struct {
	User *User  `json:"user"`
	JWT  string `json:"user_jwt"`
}

func GetCreds(f fetch.Fetcher) (*Creds, error) {
	url := "https://online-go.com/api/v1/ui/config"
	data, err := f.Fetch(url)
	if err != nil {
		return nil, err
	}
	resp := &Creds{}
	err = json.Unmarshal([]byte(data), resp)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return resp, nil
}

type OGSConnector struct {
	Creds   *Creds
	Socket  *websocket.Conn
	Room    Room
	First   int
	Exit    bool
	fetcher fetch.Fetcher
}

func NewOGSConnector(room Room, f fetch.Fetcher) (*OGSConnector, error) {
	creds, err := GetCreds(f)
	_ = creds
	if err != nil {
		return nil, err
	}

	ws, err := websocket.Dial("wss://online-go.com/socket", "", "http://localhost")
	if err != nil {
		return nil, err
	}

	return &OGSConnector{
		Creds:   creds,
		Socket:  ws,
		Room:    room,
		Exit:    false,
		fetcher: f,
	}, nil
}

func (o *OGSConnector) Send(topic string, payload map[string]any) error {
	arr := []any{topic, payload}
	data, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	_, err = o.Socket.Write(data)
	return err
}

func (o *OGSConnector) Connect(gameID int, ogsType string) error {
	payload := make(map[string]any)
	payload["player_id"] = o.Creds.User.ID
	payload["chat"] = false
	if ogsType == "game" {
		payload["game_id"] = gameID
		return o.Send("game/connect", payload)
	}
	payload["review_id"] = gameID
	return o.Send("review/connect", payload)
}

func (o *OGSConnector) ChatConnect() error {
	payload := make(map[string]any)
	payload["player_id"] = o.Creds.User.ID
	payload["username"] = o.Creds.User.Username
	payload["auth"] = o.Creds.JWT
	return o.Send("chat/connect", payload)
}

func ReadFrame(socketchan chan byte) ([]byte, error) {
	data := []byte{}
	started := false
	depth := 0
	for {
		// when the websocket is closed, ok = false
		b, ok := <-socketchan
		if !ok {
			// socket closed by reader
			return nil, nil
		}
		if !started {
			if b != '[' {
				return nil, fmt.Errorf("invalid starting byte: %c", b)
			}
			depth++
			data = append(data, b)
			started = true
		} else {
			switch b {
			case '[':
				depth++
			case ']':
				depth--
			}
			data = append(data, b)
			if depth == 0 && b == ']' {
				return data, nil
			}
		}
	}
}

func (o *OGSConnector) ReadSocketToChan(socketchan chan byte) error {
	defer close(socketchan)
	for {
		data := make([]byte, 256)
		n, err := o.Socket.Read(data)
		if err != nil {
			// this will cause the websocket to close
			// therefore ReadFrame will naturally come to an end
			return err
		}
		for _, b := range data[:n] {
			socketchan <- b
		}
		if o.Exit {
			break
		}
	}
	return nil
}

func (o *OGSConnector) Start(args map[string]any) {
	id := args["id"].(int)
	ogsType := args["ogsType"].(string)

	go o.Loop(id, ogsType) //nolint:errcheck
}

func (o *OGSConnector) End() {
	o.Exit = true
}

func (o *OGSConnector) Ping() {
	for !o.Exit {
		//30 seconds seemed just a little too long was causing connection issues
		time.Sleep(25 * time.Second)
		payload := make(map[string]any)
		payload["client"] = time.Now().UnixMilli()
		if err := o.Send("net/ping", payload); err != nil {
			log.Println(err)
			o.End()
			return
		}
	}
}

func (o *OGSConnector) Loop(gameID int, ogsType string) error {
	err := o.ChatConnect()
	if err != nil {
		return err
	}
	err = o.Connect(gameID, ogsType)
	if err != nil {
		return err
	}

	socketchan := make(chan byte)

	go o.Ping()
	go o.ReadSocketToChan(socketchan) //nolint: errcheck

	defer o.End()

	for !o.Exit {
		data, err := ReadFrame(socketchan)

		// break on error
		if err != nil {
			log.Println(err)
			break
		}

		// if err == nil and data == nil
		// then break
		if data == nil {
			log.Println("socket closed")
			break
		}

		arr := make([]any, 2)
		err = json.Unmarshal(data, &arr)
		if err != nil {
			log.Println(err)
			continue
		}
		topic := arr[0].(string)

		if topic == fmt.Sprintf("game/%d/move", gameID) {
			payload := arr[1].(map[string]any)
			move := payload["move"].([]any)

			x := int(move[0].(float64))
			y := int(move[1].(float64))

			col := core.Black
			//curColor := o.Room.State.Head.Color
			curColor := o.Room.HeadColor()
			if curColor == core.Black {
				col = core.White
			}
			//o.Room.State.PushHead(x, y, col)
			o.Room.PushHead(x, y, col)

			//frame := o.Room.State.GenerateFullFrame(core.Full)
			frame := o.Room.GenerateFullFrame(core.Full)
			evt := event.FrameEvent(frame)
			o.Room.Broadcast(evt)

		} else if topic == fmt.Sprintf("game/%d/gamedata", gameID) {
			payload := arr[1].(map[string]any)
			if _, ok := payload["winner"]; ok {
				// the game is over
				break
			}
			sgf := o.GamedataToSGF(payload)
			evt := o.Room.UploadSGF(sgf)
			o.Room.Broadcast(evt)
		} else if topic == fmt.Sprintf("review/%d/full_state", gameID) {
			/*
				nodes := arr[1].([]any)
				for _, node := range nodes {
					log.Println(node)
				}
			*/

			// eventually we can pull height, game_name, player names, etc
		} else if topic == fmt.Sprintf("review/%d/r", gameID) {
			log.Printf("review/%d/r", gameID)
			payload := arr[1].(map[string]any)
			if _, ok := payload["m"]; !ok {
				continue
			}
			moves := payload["m"].(string)

			movesArr := []*core.Stone{}
			currentColor := core.Black
			if o.First == 1 {
				currentColor = core.White
			}

			for i := 0; i < len(moves); i += 2 {
				if i+1 < len(moves) {
					coordStr := moves[i : i+2]

					switch coordStr {
					case "!1":
						//Force next move black
						currentColor = core.Black
					case "!2":
						//Force next move white
						currentColor = core.White
					case "..":
						//Pass
						movesArr = append(movesArr, &core.Stone{Coord: nil, Color: currentColor})
						currentColor = core.Opposite(currentColor)
					default:
						coord := core.LettersToCoord(coordStr)
						movesArr = append(movesArr, &core.Stone{Coord: coord, Color: currentColor})
						currentColor = core.Opposite(currentColor)
					}
				}
			}
			o.Room.AddStones(movesArr)

			// Send full board update after adding pattern
			//frame := o.Room.State.GenerateFullFrame(core.Full)
			frame := o.Room.GenerateFullFrame(core.Full)
			evt := event.FrameEvent(frame)
			o.Room.Broadcast(evt)
		}
	}
	return nil
}

func (o *OGSConnector) GamedataToSGF(gamedata map[string]any) string {
	sgf := o.GameInfoToSGF(gamedata)
	sgf += o.initStateToSGF(gamedata)

	for index, m := range gamedata["moves"].([]any) {
		arr := m.([]any)
		c := &core.Coord{X: int(arr[0].(float64)), Y: int(arr[1].(float64))}

		col := "B"

		if (index%2 == 1 && o.First == 0) || (index%2 == 0 && o.First == 1) {
			col = "W"
		}

		sgf += fmt.Sprintf(";%s[%s]", col, c.ToLetters())
	}
	sgf += ")"

	return sgf
}

func makeRank(r float64) string {
	log.Println(r)
	if r < 30 {
		return fmt.Sprintf("%dk", int(30-r+1))
	} else {
		return fmt.Sprintf("%dd", int(r-30+1))
	}
}

func (o *OGSConnector) GameInfoToSGF(gamedata map[string]any) string {
	sgf := ""

	size := int(gamedata["width"].(float64))
	komi := gamedata["komi"].(float64)
	name := gamedata["game_name"].(string)
	rules := gamedata["rules"].(string)

	players := gamedata["players"].(map[string]any)
	blackPlayer := players["black"].(map[string]any)
	whitePlayer := players["white"].(map[string]any)
	black := blackPlayer["username"].(string)
	white := whitePlayer["username"].(string)
	blackRank := makeRank(blackPlayer["rank"].(float64))
	whiteRank := makeRank(whitePlayer["rank"].(float64))
	sgf = fmt.Sprintf(
		"(;GM[1]FF[4]CA[UTF-8]SZ[%d]PB[%s]PW[%s]BR[%s]WR[%s]RU[%s]KM[%f]GN[%s]",
		size, black, white, blackRank, whiteRank, rules, komi, name)
	return sgf
}

func (o *OGSConnector) initStateToSGF(gamedata map[string]any) string {
	sgf := ""

	ip := gamedata["initial_player"].(string)

	if ip == "black" {
		o.First = 0
	} else {
		o.First = 1
	}
	initState := gamedata["initial_state"].(map[string]any)

	bstate := initState["black"].(string)
	wstate := initState["white"].(string)

	if len(bstate) > 0 {
		sgf += "AB"
		for i := 0; i < len(bstate)/2; i++ {
			sgf += fmt.Sprintf("[%s]", bstate[2*i:2*i+2])
		}
	}

	if len(wstate) > 0 {
		sgf += "AW"
		for i := 0; i < len(wstate)/2; i++ {
			sgf += fmt.Sprintf("[%s]", wstate[2*i:2*i+2])
		}
	}
	return sgf
}
