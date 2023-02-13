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
	"crypto/tls"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/francoispqt/onelog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "uptime"
)

var (
	isUp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "is_up"),
		"Is it up?",
		[]string{"endpoint_name", "endpoint_address"},
		nil,
	)

	latency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "latency_seconds"),
		"Measured latency on last scrape",
		[]string{"endpoint_name", "endpoint_address"},
		nil,
	)
)

type Exporter struct {
	Endpoints []PreprocessedEndpoint
	Logger    *onelog.Logger
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent. The sent descriptors fulfill the
// consistency and uniqueness requirements described in the Desc
// documentation.
//
// It is valid if one and the same Collector sends duplicate
// descriptors. Those duplicates are simply ignored. However, two
// different Collectors must not send duplicate descriptors.
//
// Sending no descriptor at all marks the Collector as “unchecked”,
// i.e. no checks will be performed at registration time, and the
// Collector may yield any Metric it sees fit in its Collect method.
//
// This method idempotently sends the same descriptors throughout the
// lifetime of the Collector. It may be called concurrently and
// therefore must be implemented in a concurrency safe way.
//
// If a Collector encounters an error while executing this method, it
// must send an invalid descriptor (created with NewInvalidDesc) to
// signal the error to the registry.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- isUp
	ch <- latency
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent. The
// descriptor of each sent metric is one of those returned by Describe
// (unless the Collector is unchecked, see above). Returned metrics that
// share the same descriptor must differ in their variable label
// values.
//
// This method may be called concurrently and must therefore be
// implemented in a concurrency safe way. Blocking occurs at the expense
// of total performance of rendering all registered metrics. Ideally,
// Collector implementations support concurrent readers.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(len(e.Endpoints))

	for _, endpoint := range e.Endpoints {
		go func(endpoint PreprocessedEndpoint) {
			defer wg.Done()

			// Create HTTP calls
			request, err := http.NewRequestWithContext(ctx, endpoint.Method, endpoint.Address, nil)
			if err != nil {
				e.Logger.Error(err.Error())
				return
			}

			httpClient := http.Client{
				Timeout: time.Duration(uint64(endpoint.Timeout) * 1e+9),
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						Certificates:       endpoint.TLSConfiguration.Certificates,
						RootCAs:            endpoint.TLSConfiguration.RootCA,
						InsecureSkipVerify: endpoint.TLSConfiguration.InsecureSkipVerify,
					},
				},
			}

			response, err := httpClient.Do(request)
			if err != nil {
				e.Logger.Error(err.Error())
				return
			}

			var isUpValue bool = false
			if !checkStatusCode(response.StatusCode, endpoint.SuccessfulStatusCode) {
				isUpValue = true
			}

			if endpoint.InverseStatus {
				isUpValue = !isUpValue
			}

			var isUpResult float64 = 0
			if isUpValue {
				isUpResult = 1
			}

			ch <- prometheus.MustNewConstMetric(
				isUp,
				prometheus.GaugeValue,
				isUpResult,
				endpoint.Name,
				endpoint.Address,
			)
			ch <- prometheus.MustNewConstMetric(
				latency,
				prometheus.GaugeValue,
				time.Since(start).Seconds(),
				endpoint.Name,
				endpoint.Address,
			)

		}(endpoint)
	}

	wg.Wait()
}

// Return true if the statusCode matches what's expected.
func checkStatusCode(statusCode int, expected string) bool {
	strStatusCode := strconv.Itoa(statusCode)
	// quick check
	if strStatusCode == expected {
		return true
	}

	// if it's not, we parse the expected string, if it contains any 'x' characters
	if len(expected) != len(strStatusCode) {
		// as this is pretty much impossible to check
		return false
	}

	var ok = true
	for i := 0; i < len(strStatusCode); i++ {
		if string(expected[i]) == "x" {
			continue
		}

		if string(expected[i]) == string(strStatusCode[i]) {
			continue
		} else {
			ok = false
			break
		}
	}

	return ok
}
