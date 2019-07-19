// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package gcp

import (
	"context"

	"cloud.google.com/go/pubsub"

	"github.com/elastic/beats/libbeat/cmd/instance"
	"github.com/elastic/beats/x-pack/functionbeat/manager/beater"
	"github.com/elastic/beats/x-pack/functionbeat/provider/gcp/gcp/transformer"
	_ "github.com/elastic/beats/x-pack/functionbeat/provider/gcp/include"
)

func ForwardPubSubEvent(ctx context.Context, m pubsub.Message) {
	fb, err := initFunctionbeat()
	if err != nil {
		fb.log.Debugf("Failed to init functionbeat: %+v", err)
		return
	}

	event, err := transformer.PubSub(ctx, m)
	if err != nil {
		fb.log.Debugf("Cannot transform event: %+v", err)
		return
	}

	fb.log.Debugf("The handler received Pub/Sub event: %+v", event)

	if err := client.Publish(event); err != nil {
		c.log.Errorf("Could not publish event to the pipeline, error: %+v", err)
		return err
	}
	client.Wait()
}

func initFunctionbeat() (*Functionbeat, error) {
	b, err := instance.NewInitializedBeat(instance.Settings{Name: "functionbeat"})
	if err != nil {
		return nil, err
	}

	fn, err := beater.New(b, b.Config)
	if err != nil {
		return nil, err
	}

	return fn, nil
}
