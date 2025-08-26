/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jarednogo/board/backend/loader"
	"github.com/jarednogo/board/backend/server"
	"golang.org/x/net/websocket"
)

func Serve(url string) {
	// wrapping serve in the log.Fatal call ensures
	// that when it's called in a goroutine and there's an error
	// we end the program and print the error
	log.Fatal(http.ListenAndServe(url, nil))
}

func main() {
	// create dirs
	loader.Setup()

	// create empty config
	cfg := websocket.Config{}

	// create new server, load rooms, defer save
	s := server.NewServer()
	s.Load()
	defer s.Save()

	// create new websocket server
	ws := websocket.Server{
		Config:    cfg,
		Handshake: nil,
		Handler:   s.Handler,
	}
	http.Handle("/", ws)

	port := "9000"
	host := "localhost"
	url := fmt.Sprintf("%s:%s", host, port)

	log.Println("Listening on", url)

	// get ready to catch signals
	cancelChan := make(chan os.Signal, 1)

	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	// start message loop
	go s.MessageLoop()

	// start http loop
	go Serve(url)

	// catch cancel signal
	sig := <-cancelChan

	log.Printf("Caught signal %v", sig)
	log.Println("Shutting down gracefully")
}
