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

package input_logfile

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	input "github.com/elastic/beats/v7/filebeat/input/v2"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/go-concert/ctxtool"
	"github.com/elastic/go-concert/unison"
)

type Harvester interface {
	Name() string
	Test(Source, input.TestContext) error
	Run(input.Context, Source, Cursor, Publisher) error
}

type HarvesterGroup struct {
	manager      *InputManager
	readers      map[string]context.CancelFunc
	pipeline     beat.PipelineConnector
	harvester    Harvester
	cleanTimeout time.Duration
	store        *store
	tg           unison.TaskGroup
}

func (hg *HarvesterGroup) Run(ctx input.Context, s Source) error {
	log := ctx.Logger.With("source", s.Name())
	log.Debug("Starting harvester for file")

	harvesterCtx, cancelHarvester := context.WithCancel(ctxtool.FromCanceller(ctx.Cancelation))
	ctx.Cancelation = harvesterCtx

	resource, err := hg.manager.lock(ctx, s.Name())
	if err != nil {
		return err
	}

	if _, ok := hg.readers[s.Name()]; ok {
		log.Debug("A harvester is already running for file")
		return nil
	}
	hg.readers[s.Name()] = cancelHarvester

	hg.store.UpdateTTL(resource, hg.cleanTimeout)

	client, err := hg.pipeline.ConnectWith(beat.ClientConfig{
		CloseRef:   ctx.Cancelation,
		ACKHandler: newInputACKHandler(ctx.Logger),
	})
	if err != nil {
		return err
	}

	cursor := makeCursor(hg.store, resource)
	publisher := &cursorPublisher{canceler: ctx.Cancelation, client: client, cursor: &cursor}

	go func() {
		defer client.Close()
		defer log.Debug("Stopped harvester for file")
		defer cancelHarvester()
		defer releaseResource(resource)

		defer func() {
			if v := recover(); v != nil {
				err := fmt.Errorf("harvester panic with: %+v\n%s", v, debug.Stack())
				ctx.Logger.Errorf("Harvester crashed with: %+v", err)
			}
		}()

		err := hg.harvester.Run(ctx, s, cursor, publisher)
		if err != nil {
			log.Errorf("Harvester stopped: %v", err)
		}
	}()
	return nil
}

func (hg *HarvesterGroup) Cancel(s Source) error {
	if cancel, ok := hg.readers[s.Name()]; ok {
		cancel()
		return nil
	}
	return fmt.Errorf("no such harvester %s", s.Name())
}

func releaseResource(resource *resource) {
	resource.lock.Unlock()
	resource.Release()
}
