/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type client struct {
	socket *websocket.Conn
	room   string
}

func NewClient(serverAddr, room string) (*client, error) {
	tls := !strings.HasPrefix(serverAddr, "localhost")
	path := fmt.Sprintf("/socket/b/%s", room)
	scheme := "http"
	wsScheme := "ws"
	if tls {
		scheme += "s"
		wsScheme += "s"
	}
	serverURL := url.URL{Scheme: wsScheme, Host: serverAddr, Path: path}

	cfg, _ := websocket.NewConfig(serverURL.String(), scheme+"://"+serverAddr)
	// cloudfront has a default aws rule for blocking connections with no user agent
	cfg.Header.Set("User-Agent", "load-testing-client")

	ws, err := websocket.DialConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &client{socket: ws, room: room}, nil
}

func (c *client) send(msgStr string) error {
	msg := []byte(msgStr)
	// allocate 4 bytes
	buf := make([]byte, 4)

	// encode little-endian
	binary.LittleEndian.PutUint32(buf, uint32(len(msg)))
	_, err := c.socket.Write(buf)
	if err != nil {
		return err
	}
	_, err = c.socket.Write(msg)
	return err
}

func (c *client) receive() bool {
	// read in 4 bytes (length of rest of message)
	lengthArray := make([]byte, 4)
	_, err := c.socket.Read(lengthArray)
	return err != nil
}

func (c *client) Close() error {
	return c.socket.Close()
}

func main() {
	numRoomsp := flag.Int("r", 5, "number of rooms to connect to")
	numClientsp := flag.Int("c", 10, "number of clients to create per room")
	serverAddrp := flag.String("a", "localhost:8080", "server address")
	flag.Parse()
	serverAddr := ""
	if *serverAddrp != "" {
		serverAddr = *serverAddrp
	}
	var numRooms int
	var numClients int
	if *numRoomsp != 0 {
		numRooms = *numRoomsp
	}
	if *numClientsp != 0 {
		numClients = *numClientsp
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	latencies := []time.Duration{}
	for r := 0; r < numRooms; r++ {
		room := fmt.Sprintf("room%d", r)
		for i := 0; i < numClients; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				start := time.Now()
				c, err := NewClient(serverAddr, room)
				if err != nil {
					return
				}

				defer c.Close() //nolint:errcheck
				c.receive()
				payload := fmt.Sprintf(
					`{"event": "comment", "value": "room%d-client%d"}`,
					r,
					i,
				)
				err = c.send(payload)
				if err != nil {
					return
				}
				latency := time.Since(start)
				mu.Lock()
				latencies = append(latencies, latency)
				mu.Unlock()
			}()
		}
	}
	wg.Wait()

	report(latencies, numRooms, numClients)
}

func report(latencies []time.Duration, numRooms, numClients int) {
	if len(latencies) == 0 {
		fmt.Println("No successful clients connected")
		return
	}

	var totalLatency time.Duration
	minLatency := latencies[0]
	maxLatency := latencies[0]

	for _, l := range latencies {
		totalLatency += l
		if l < minLatency {
			minLatency = l
		}
		if l > maxLatency {
			maxLatency = l
		}
	}

	avgLatency := totalLatency / time.Duration(len(latencies))

	fmt.Println("=== Load Test Report ===")
	fmt.Printf("Number of rooms: %d\n", numRooms)
	fmt.Printf("Clients per room: %d\n", numClients)
	fmt.Printf("Total clients: %d\n", numClients*numRooms)
	fmt.Printf("Latency (min/avg/max): %v / %v / %v\n", minLatency, avgLatency, maxLatency)
}
