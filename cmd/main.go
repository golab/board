/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golab/board/pkg/app"
	"github.com/golab/board/pkg/config"
	"github.com/golab/board/pkg/logx"
)

var version = "dev"

func main() {
	// set up root logger
	level := logx.LogLevelInfo
	logger := logx.NewDefaultLogger(level)

	// read in config
	cfgFile := flag.String("f", "", "Path to config file")
	flag.Parse()
	var cfg *config.Config

	if *cfgFile == "" {
		cfg = config.Default()
	} else if loadedCfg, err := config.New(*cfgFile); err != nil {
		logger.Error("failed loading config", "file", *cfgFile, "err", err)
		cfg = config.Default()
	} else {
		logger.Info("loaded config", "file", *cfgFile)
		cfg = loadedCfg
	}

	cfg.Version = version

	safe := *cfg
	safe.Redact()
	logger.Info("running config", "config", fmt.Sprintf("%v", safe))

	// setup routes
	a, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("error in setup", "err", err)
		return
	}

	// start everything
	a.Hub.Load()
	defer a.Hub.Save()
	host := cfg.Server.Host
	port := cfg.Server.Port
	url := fmt.Sprintf("%s:%d", host, port)

	// get ready to catch signals
	cancelChan := make(chan os.Signal, 1)

	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	// run server
	go func() {
		logger.Info("listening on", "url", url)
		err = http.ListenAndServe(url, a.Router)
		logger.Error("error listening", "err", err)
	}()

	// catch cancel signal
	sig := <-cancelChan

	logger.Info("caught signal", "signal", fmt.Sprintf("%v", sig))
	logger.Info("shutting down gracefully")
}
