// MIT License
//
// Copyright (c) 2023 Reinaldy Rafli <aldy505@proton.me>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/francoispqt/onelog"
	"github.com/jinzhu/configor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	configFile := flag.String("config-file", "", "Path to configuration file")
	listenAddress := flag.String("web.listen-address", ":9428", "Address on which to expose metrics and web interface")
	metricsPath := flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")

	flag.Parse()

	// Read configuration from file
	var config Configuration
	err := configor.Load(&config, *configFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate the logger instance. It should be zero-allocated, so we'll use onelog
	logger := onelog.New(os.Stdout, onelog.ALL)

	exporter := &Exporter{
		Endpoints: config.Endpoints,
		Logger:    logger,
	}

	r := prometheus.NewRegistry()
	r.MustRegister(exporter)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`<html>
             <head><title>Uptime Exporter</title></head>
             <body>
             <h1>Uptime Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             <p><a href='/health'>Health</a></p>
             </body>
             </html>`))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.Handle(*metricsPath, promhttp.HandlerFor(r, promhttp.HandlerOpts{
		MaxRequestsInFlight: 1,
		Timeout:             time.Minute * 5,
	}))

	server := &http.Server{
		Addr:              *listenAddress,
		Handler:           mux,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		log.Printf("Server is starting on %s", *listenAddress)

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listening to HTTP: %v", err)
		}
	}()

	<-exitSignal

	log.Printf("Recevied exit signal")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutting down server: %v", err)
	}
}
