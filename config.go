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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type Configuration struct {
	Endpoints []Endpoint `required:"true"`
}

type Endpoint struct {
	Name                 string           `required:"true"`
	Address              string           `required:"true"`
	Method               string           `required:"false" default:"GET"`
	Timeout              uint16           `required:"false" default:"30"`
	SuccessfulStatusCode string           `required:"false" default:"2xx" yaml:"successful_status_code" json:"successful_status_code" toml:"successful_status_code"`
	InverseStatus        bool             `required:"false" default:"false" yaml:"inverse_status" json:"inverse_status" toml:"inverse_status"`
	TLSConfiguration     TLSConfiguration `required:"false" yaml:"tls_configuration" json:"tls_configuration" toml:"tls_configuration"`
}

type TLSConfiguration struct {
	CertificateAuthorityPath string `required:"false" yaml:"certificate_authority_path" json:"certificate_authority_path" toml:"certificate_authority_path"`
	ClientCertificatePath    string `required:"false" yaml:"client_certificate_path" json:"client_certificate_path" toml:"client_certificate_path"`
	ClientKeyPath            string `required:"false" yaml:"client_key_path" json:"client_key_path" toml:"client_key_path"`
	InsecureSkipVerify       bool   `required:"false" default:"false" yaml:"insecure_skip_verify" json:"insecure_skip_verify" toml:"insecure_skip_verify"`
}

type PreprocessedEndpoint struct {
	Name                 string
	Address              string
	Method               string
	Timeout              uint16
	SuccessfulStatusCode string
	InverseStatus        bool
	TLSConfiguration     PreproccessedTLSConfiguration
}

type PreproccessedTLSConfiguration struct {
	Certificates       []tls.Certificate
	RootCA             *x509.CertPool
	InsecureSkipVerify bool
}

func processConfiguration(endpoints []Endpoint) ([]PreprocessedEndpoint, error) {
	var preprocessedEndpoints []PreprocessedEndpoint

	for _, endpoint := range endpoints {
		var certificates []tls.Certificate = nil
		var rootCA *x509.CertPool = nil

		if endpoint.TLSConfiguration.CertificateAuthorityPath != "" {
			// Read CA file
			file, err := os.ReadFile(endpoint.TLSConfiguration.CertificateAuthorityPath)
			if err != nil {
				return nil, fmt.Errorf("reading certificate authority file: %s", err.Error())
			}

			rootCA = x509.NewCertPool()
			ok := rootCA.AppendCertsFromPEM(file)
			if !ok {
				return nil, fmt.Errorf("invalid pem format for certificate authority")
			}
		}

		if endpoint.TLSConfiguration.ClientCertificatePath != "" && endpoint.TLSConfiguration.ClientKeyPath != "" {
			// Read certificate file
			certificateFile, err := os.ReadFile(endpoint.TLSConfiguration.ClientCertificatePath)
			if err != nil {
				return nil, fmt.Errorf("reading certificate file: %s", err.Error())
			}

			keyFile, err := os.ReadFile(endpoint.TLSConfiguration.ClientKeyPath)
			if err != nil {
				return nil, fmt.Errorf("reading key file: %s", err.Error())
			}

			certificate, err := tls.X509KeyPair(certificateFile, keyFile)
			if err != nil {
				return nil, fmt.Errorf("creating x509 key pair: %s", err.Error())
			}

			certificates = append(certificates, certificate)
		}

		preprocessedEndpoints = append(preprocessedEndpoints, PreprocessedEndpoint{
			Name:                 endpoint.Name,
			Address:              endpoint.Address,
			Method:               endpoint.Method,
			Timeout:              endpoint.Timeout,
			SuccessfulStatusCode: endpoint.SuccessfulStatusCode,
			InverseStatus:        endpoint.InverseStatus,
			TLSConfiguration: PreproccessedTLSConfiguration{
				Certificates:       certificates,
				RootCA:             rootCA,
				InsecureSkipVerify: endpoint.TLSConfiguration.InsecureSkipVerify,
			},
		})
	}

	return preprocessedEndpoints, nil
}
