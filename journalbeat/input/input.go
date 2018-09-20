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
	"sync"

	"github.com/elastic/beats/journalbeat/checkpoint"
	"github.com/elastic/beats/journalbeat/reader"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

// Input manages readers and forwards entries from journals.
type Input struct {
	readers  []*reader.Reader
	done     chan struct{}
	config   Config
	pipeline beat.Pipeline
	states   map[string]checkpoint.JournalState
}

// New returns a new Inout
func New(
	c *common.Config,
	pipeline beat.Pipeline,
	done chan struct{},
	states map[string]checkpoint.JournalState,
) (*Input, error) {
	config := DefaultConfig
	if err := c.Unpack(&config); err != nil {
		return nil, err
	}
	var readers []*reader.Reader
	if len(config.Paths) == 0 {
		cfg := reader.Config{
			Path:          reader.LocalSystemJournalID, // used to identify the state in the registry
			Backoff:       config.Backoff,
			MaxBackoff:    config.MaxBackoff,
			BackoffFactor: config.BackoffFactor,
			Seek:          config.Seek,
		}

		state := states[reader.LocalSystemJournalID]
		r, err := reader.NewLocal(cfg, done, state)
		if err != nil {
			return nil, fmt.Errorf("error creating reader for local journal: %v", err)
		}
		readers = append(readers, r)
	}

	for _, p := range config.Paths {
		cfg := reader.Config{
			Path:          p,
			Backoff:       config.Backoff,
			MaxBackoff:    config.MaxBackoff,
			BackoffFactor: config.BackoffFactor,
			Seek:          config.Seek,
		}
		state := states[p]
		r, err := reader.New(cfg, done, state)
		if err != nil {
			return nil, fmt.Errorf("error creating reader for journal: %v", err)
		}
		readers = append(readers, r)
	}

	return &Input{
		readers:  readers,
		done:     done,
		config:   config,
		pipeline: pipeline,
		states:   states,
	}, nil
}

// Run connects to the output, collects entries from the readers
// and then publishes the events.
func (i *Input) Run() {
	client, err := i.pipeline.ConnectWith(beat.ClientConfig{
		PublishMode:   beat.GuaranteedSend,
		EventMetadata: common.EventMetadata{},
		Meta:          nil,
		Processor:     nil,
		ACKCount: func(n int) {
			logp.Info("journalbeat successfully published %d events", n)
		},
	})
	if err != nil {
		logp.Err("Error connecting: %v", err)
		return
	}
	defer client.Close()

	for {
		select {
		case <-i.done:
			return
		default:
			i.publishAll(client)
		}
	}

}

func (i *Input) publishAll(client beat.Client) {
	out := make(chan *beat.Event)
	var wg sync.WaitGroup

	merge := func(in chan *beat.Event) {
		wg.Add(1)

		go func(c chan *beat.Event) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(in)
	}
	go func() {
		wg.Wait()
		close(out)
	}()

	for _, r := range i.readers {
		c := r.Follow()
		merge(c)
	}

	for e := range out {
		client.Publish(*e)
	}
}

// Stop stops all readers of the input.
func (i *Input) Stop() {
	for _, r := range i.readers {
		r.Close()
	}
}

// Wait waits until all readers are done.
func (i *Input) Wait() {
	i.Stop()
}
