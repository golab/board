/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package socket

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"log"

	"github.com/jarednogo/board/pkg/core"
	"golang.org/x/net/websocket"
)

func SendEvent(conn *websocket.Conn, evt *core.EventJSON) {
	// marshal event back into data
	data, err := json.Marshal(evt)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func ReceiveEvent(ws *websocket.Conn) (*core.EventJSON, error) {
	data, err := ReadPacket(ws)
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

func ReadPacket(ws *websocket.Conn) ([]byte, error) {
	// read in 4 bytes (length of rest of message)
	lengthArray := make([]byte, 4)
	_, err := ws.Read(lengthArray)
	if err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lengthArray)

	// read in the rest of the data
	var data []byte

	if length > 1024 {
		data, err = ReadBytes(ws, int(length))
		if err != nil {
			return nil, err
		}
	} else {
		data = make([]byte, length)
		_, err := ws.Read(data)

		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func ReadBytes(ws *websocket.Conn, size int) ([]byte, error) {
	chunkSize := 64
	message := []byte{}
	for len(message) < size {
		l := size - len(message)
		if l > chunkSize {
			l = chunkSize
		}
		temp := make([]byte, l)
		n, err := ws.Read(temp)
		if err != nil {
			return nil, err
		}
		message = append(message, temp[:n]...)
	}

	return message, nil

}

func EncodeSend(ws *websocket.Conn, data string) {
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	length := uint32(len(encoded))
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, length)

	_, err := ws.Write(buf)
	if err != nil {
		log.Println(err)
	}
	_, err = ws.Write([]byte(encoded))
	if err != nil {
		log.Println(err)
	}
}
