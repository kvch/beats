// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package tlscommon

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/elastic/beats/v7/libbeat/logp"
)

// TLSConfig is the interface used to configure a tcp client or server from a `Config`
type TLSConfig struct {

	// List of allowed SSL/TLS protocol versions. Connections might be dropped
	// after handshake succeeded, if TLS version in use is not listed.
	Versions []TLSVersion

	// Configure SSL/TLS verification mode used during handshake. By default
	// VerifyFull will be used.
	Verification TLSVerificationMode

	// List of certificate chains to present to the other side of the
	// connection.
	Certificates []tls.Certificate

	// Set of root certificate authorities use to verify server certificates.
	// If RootCAs is nil, TLS might use the system its root CA set (not supported
	// on MS Windows).
	RootCAs *x509.CertPool

	// Set of root certificate authorities use to verify client certificates.
	// If ClientCAs is nil, TLS might use the system its root CA set (not supported
	// on MS Windows).
	ClientCAs *x509.CertPool

	// List of supported cipher suites. If nil, a default list provided by the
	// implementation will be used.
	CipherSuites []uint16

	// Types of elliptic curves that will be used in an ECDHE handshake. If empty,
	// the implementation will choose a default.
	CurvePreferences []tls.CurveID

	// Renegotiation controls what types of renegotiation are supported.
	// The default, never, is correct for the vast majority of applications.
	Renegotiation tls.RenegotiationSupport

	// ClientAuth controls how we want to verify certificate from a client, `none`, `optional` and
	// `required`, default to required. Do not affect TCP client.
	ClientAuth tls.ClientAuthType

	// CASha256 is the CA certificate pin, this is used to validate the CA that will be used to trust
	// the server certificate.
	CASha256 []string

	// time returns the current time as the number of seconds since the epoch.
	// If time is nil, TLS uses time.Now.
	time func() time.Time
}

// ToConfig generates a tls.Config object. Note, you must use BuildModuleConfig to generate a config with
// ServerName set, use that method for servers with SNI.
func (c *TLSConfig) ToConfig() *tls.Config {
	if c == nil {
		return &tls.Config{}
	}

	minVersion, maxVersion := extractMinMaxVersion(c.Versions)

	// When we are using the CAsha256 pin to validate the CA used to validate the chain,
	// or when we are using 'certificate' TLS verification mode, we add a custom callback
	verifyConnectionFn := makeVerifyConnection(c)

	fmt.Println(c.Verification)
	insecure := c.Verification != VerifyStrict
	if c.Verification == VerifyNone {
		logp.NewLogger("tls").Warn("SSL/TLS verifications disabled.")
	}
	return &tls.Config{
		MinVersion:         minVersion,
		MaxVersion:         maxVersion,
		Certificates:       c.Certificates,
		RootCAs:            c.RootCAs,
		ClientCAs:          c.ClientCAs,
		InsecureSkipVerify: insecure,
		CipherSuites:       c.CipherSuites,
		CurvePreferences:   c.CurvePreferences,
		Renegotiation:      c.Renegotiation,
		ClientAuth:         c.ClientAuth,
		Time:               c.time,
		VerifyConnection:   verifyConnectionFn,
	}
}

// BuildModuleConfig takes the TLSConfig and transform it into a `tls.Config`.
func (c *TLSConfig) BuildModuleConfig(host string) *tls.Config {
	if c == nil {
		// use default TLS settings, if config is empty.
		return &tls.Config{ServerName: host}
	}

	config := c.ToConfig()
	config.ServerName = host
	return config
}

func makeVerifyConnection(cfg *TLSConfig) func(tls.ConnectionState) error {
	pin := len(cfg.CASha256) > 0

	switch cfg.Verification {
	case VerifyFull:
		return func(cs tls.ConnectionState) error {
			dnsnames := cs.PeerCertificates[0].DNSNames
			var serverName string
			if len(dnsnames) == 0 || len(dnsnames) == 1 && dnsnames[0] == "" {
				serverName = cs.PeerCertificates[0].Subject.CommonName
			} else {
				serverName = dnsnames[0]
			}
			if len(serverName) > 0 && len(cs.ServerName) > 0 && serverName != cs.ServerName {
				return x509.HostnameError{Certificate: cs.PeerCertificates[0], Host: cs.ServerName}
			}
			opts := x509.VerifyOptions{
				Roots:         cfg.RootCAs,
				Intermediates: x509.NewCertPool(),
			}
			for _, cert := range cs.PeerCertificates[1:] {
				opts.Intermediates.AddCert(cert)
			}
			_, err := cs.PeerCertificates[0].Verify(opts)
			if err != nil {
				return err
			}

			if pin {
				verifiedChains := [][]*x509.Certificate{cs.PeerCertificates}
				return verifyCAPin(cfg.CASha256, verifiedChains)
			}
			return nil
		}
	case VerifyCertificate:
		return func(cs tls.ConnectionState) error {
			opts := x509.VerifyOptions{
				Roots:         cfg.RootCAs,
				Intermediates: x509.NewCertPool(),
			}
			for _, cert := range cs.PeerCertificates[1:] {
				opts.Intermediates.AddCert(cert)
			}
			_, err := cs.PeerCertificates[0].Verify(opts)
			if err != nil {
				return err
			}

			if pin {
				verifiedChains := [][]*x509.Certificate{cs.PeerCertificates}
				return verifyCAPin(cfg.CASha256, verifiedChains)
			}
			return nil
		}
	default:
	}

	return nil

}
