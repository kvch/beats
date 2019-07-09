// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package gcp

import (
	"context"

	"cloud.google.com/go/pubsub"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/x-pack/functionbeat/function/core"
	"github.com/elastic/beats/x-pack/functionbeat/function/provider"
)

// PubSub represents a Google Cloud function which reads event from Google Pub/Sub triggers.
type PubSub struct {
	config *functionConfig
}

// NewPubSub returns a new function to read from Google Pub/Sub.
func NewPubSub(provider provider.Provider, config *common.Config) (provider.Function, error) {
	functionConfig := &functionConfig{}
	if err := cfg.Unpack(functionConfig); err != nil {
		return nil, err
	}
	return &PubSub{
		config: functionConfig,
	}, nil
}

// Run start the AWS lambda handles and will transform any events received to the pipeline.
func (c *CloudwatchLogs) Run(_ context.Context, client core.Client) error {
	return nil
}

func handlePubSubMessage(ctx context.Context, m pubsub.Message) error {

}

// Name returns the name of the function.
func (p *PubSub) Name() string {
	return p.config.Name
}
