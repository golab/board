/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jarednogo/board/pkg/config"
)

func main() {
	// read in config
	cfgFile := flag.String("f", "", "Path to config file")
	flag.Parse()
	var cfg *config.Config

	if *cfgFile == "" {
		cfg = config.Default()
	} else if loadedCfg, err := config.New(*cfgFile); err != nil {
		log.Println("failed to load config:", *cfgFile, err)
		cfg = config.Default()
	} else {
		log.Println("successfully loaded config:", *cfgFile)
		cfg = loadedCfg
	}

	log.Println("running config:", cfg)

	// setup routes
	r, err := Setup(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	// start everything
	host := cfg.Server.Host
	port := cfg.Server.Port
	url := fmt.Sprintf("%s:%d", host, port)

	// get ready to catch signals
	cancelChan := make(chan os.Signal, 1)

	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	// run server
	go func() {
		log.Println("Listening on", url)
		log.Fatal(http.ListenAndServe(url, r))
	}()

	// catch cancel signal
	sig := <-cancelChan

	log.Printf("Caught signal %v", sig)
	log.Println("Shutting down gracefully")

}
