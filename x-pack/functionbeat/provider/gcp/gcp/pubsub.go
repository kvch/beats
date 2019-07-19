// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package gcp

import (
	"context"

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
	if err := config.Unpack(functionConfig); err != nil {
		return nil, err
	}
	return &PubSub{
		config: functionConfig,
	}, nil
}

// Run start
func (c *PubSub) Run(_ context.Context, _ core.Client) error {
	return nil
}

// Name returns the name of the function.
func (p *PubSub) Name() string {
	return "pubsub"
}
