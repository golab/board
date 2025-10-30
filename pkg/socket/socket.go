/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package socket

import (
	"encoding/binary"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jarednogo/board/pkg/core"
	"golang.org/x/net/websocket"
)

type RoomConn interface {
	SendEvent(evt *core.EventJSON) error
	ReceiveEvent() (*core.EventJSON, error)
	Close() error
	ID() string
}

// WebsocketRoomConn is a thin wrapper around *websocket.Conn
type WebsocketRoomConn struct {
	ws *websocket.Conn
	id string
}

func NewWebsocketRoomConn(ws *websocket.Conn) RoomConn {
	// assign the new connection a new id
	id := uuid.New().String()

	return &WebsocketRoomConn{ws, id}
}

func (wrc *WebsocketRoomConn) ID() string {
	return wrc.id
}

func (wrc *WebsocketRoomConn) SendEvent(evt *core.EventJSON) error {
	// marshal event back into data
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, err = wrc.ws.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (wrc *WebsocketRoomConn) ReceiveEvent() (*core.EventJSON, error) {
	data, err := wrc.readPacket()
	if err != nil {
		return nil, err
	}

	// turn data into json
	evt := &core.EventJSON{}
	if err := json.Unmarshal(data, evt); err != nil {
		return nil, err
	}
	return evt, nil
}

func (wrc *WebsocketRoomConn) Close() error {
	return wrc.ws.Close()
}

func (wrc *WebsocketRoomConn) readPacket() ([]byte, error) {
	// read in 4 bytes (length of rest of message)
	lengthArray := make([]byte, 4)
	_, err := wrc.ws.Read(lengthArray)
	if err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lengthArray)

	// read in the rest of the data
	var data []byte

	if length > 1024 {
		data, err = wrc.readBytes(int(length))
		if err != nil {
			return nil, err
		}
	} else {
		data = make([]byte, length)
		_, err := wrc.ws.Read(data)

		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (wrc *WebsocketRoomConn) readBytes(size int) ([]byte, error) {
	chunkSize := 64
	message := []byte{}
	for len(message) < size {
		l := size - len(message)
		if l > chunkSize {
			l = chunkSize
		}
		temp := make([]byte, l)
		n, err := wrc.ws.Read(temp)
		if err != nil {
			return nil, err
		}
		message = append(message, temp[:n]...)
	}

	return message, nil
}
