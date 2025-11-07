/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package event

import (
	"encoding/binary"
	"encoding/json"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type EventChannel interface {
	SendEvent(evt Event) error
	ReceiveEvent() (Event, error)
	OnConnect()
	Close() error
	ID() string
}

// WebsocketEventChannel is a thin wrapper around *websocket.Conn
type WebsocketEventChannel struct {
	ws *websocket.Conn
	id string
}

func NewWebsocketEventChannel(ws *websocket.Conn) EventChannel {
	// assign the new connection a new id
	id := uuid.New().String()

	return &WebsocketEventChannel{ws, id}
}

func (ec *WebsocketEventChannel) ID() string {
	return ec.id
}

func (ec *WebsocketEventChannel) SendEvent(evt Event) error {
	// marshal event back into data
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, err = ec.ws.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (ec *WebsocketEventChannel) OnConnect() {
}

func (ec *WebsocketEventChannel) ReceiveEvent() (Event, error) {
	data, err := ec.readPacket()
	if err != nil {
		return nil, err
	}

	// turn data into json
	var evt Event
	if evt, err = EventFromJSON(data); err != nil {
		return nil, err
	}
	return evt, nil
}

func (ec *WebsocketEventChannel) Close() error {
	return ec.ws.Close()
}

func (ec *WebsocketEventChannel) readPacket() ([]byte, error) {
	// read in 4 bytes (length of rest of message)
	lengthArray := make([]byte, 4)
	_, err := ec.ws.Read(lengthArray)
	if err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lengthArray)

	// read in the rest of the data
	var data []byte

	if length > 1024 {
		data, err = ec.readBytes(int(length))
		if err != nil {
			return nil, err
		}
	} else {
		data = make([]byte, length)
		_, err := ec.ws.Read(data)

		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (ec *WebsocketEventChannel) readBytes(size int) ([]byte, error) {
	chunkSize := 64
	message := []byte{}
	for len(message) < size {
		l := size - len(message)
		if l > chunkSize {
			l = chunkSize
		}
		temp := make([]byte, l)
		n, err := ec.ws.Read(temp)
		if err != nil {
			return nil, err
		}
		message = append(message, temp[:n]...)
	}

	return message, nil
}
