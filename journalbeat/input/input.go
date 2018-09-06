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

package input

import (
	"fmt"
	"time"

	"github.com/elastic/beats/journalbeat/reader"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

type Config struct {
	Paths []string `config:"paths"`

	Backoff       time.Duration `config:"backoff" validate:"min=0,nonzero"`
	BackoffFactor int           `config:"backoff_factor" validate:"min=1"`
	MaxBackoff    time.Duration `config:"max_backoff" validate:"min=0,nonzero"`

	Matches map[string]string `config:"matches"`
	Seek    string            `config:"seek"`
}

var (
	DefaultConfig = Config{}
)

func (c *Config) Validate() error {
	correctSeek := false
	for _, s := range []string{"cursor", "head", "tail"} {
		if c.Seek == s {
			correctSeek = true
		}
	}

	if !correctSeek {
		return fmt.Errorf("incorrect value for seek: %s. possible values: cursor, head, tail", c.Seek)
	}

	return nil
}

type Input struct {
	readers []*reader.Reader
	done    chan struct{}
	config  Config
	client  beat.Client
}

func New(c *common.Config, client beat.Client, done chan struct{}) *Input {
	config := DefaultConfig
	if err := c.Unpack(&config); err != nil {
		// TODO
		return nil
	}
	var readers []*reader.Reader
	if len(config.Paths) == 0 {
		cfg := reader.Config{
			Backoff:       config.Backoff,
			MaxBackoff:    config.MaxBackoff,
			BackoffFactor: config.BackoffFactor,
		}
		r, err := reader.NewLocal(cfg, done)
		if err != nil {
			logp.Debug("input", "Error creating reader: %v", err)
			return nil
		}
		readers = append(readers, r)
	}

	for _, p := range config.Paths {
		cfg := reader.Config{
			Path:          p,
			Backoff:       config.Backoff,
			MaxBackoff:    config.MaxBackoff,
			BackoffFactor: config.BackoffFactor,
		}
		r, err := reader.New(cfg, done)
		if err != nil {
			logp.Debug("input", "Error creating reader: %v", err)
			continue
		}
		readers = append(readers, r)
	}

	return &Input{
		readers: readers,
		done:    done,
		client:  client,
		config:  config,
	}
}

func (i *Input) Run() {
	for {
		select {
		case <-i.done:
			return
		default:
			for _, r := range i.readers {
				for e := range r.Follow() {
					i.client.Publish(*e)
				}
			}
		}
	}

}

func (i *Input) Stop() {

}

func (i *Input) Wait() {

}
