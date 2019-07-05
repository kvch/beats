// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package gcp

import (
	"fmt"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/x-pack/functionbeat/function/provider"
)

const (
	// AWS lambda currently support go 1.x as a runtime.
	runtime          = "go111"        // Golang 1.11
	entryPoint       = "functionbeat" // entrypoint
	sourceArchiveURL = "gs://%s"      // location of functionbeat.zip
)

// CLIManager interacts with the AWS Lambda API to deploy, update or remove a function.
// It will take care of creating the main lambda function and ask for each function type for the
// operation that need to be executed to connect the lambda to the triggers.
type CLIManager struct {
	templateBuilder *restAPITemplateBuilder
	log             *logp.Logger
	config          *Config
	functionConfig  functionConfig
}

// Deploy delegate deploy to the actual function implementation.
func (c *CLIManager) Deploy(name string) error {
	c.log.Debugf("Deploying function: %s", name)
	defer c.log.Debugf("Deploy finish for function '%s'", name)

	// TODO deploy function

	c.log.Debugf("Successfully created function '%s'", name)
	return nil
}

// Update updates the function.
func (c *CLIManager) Update(name string) error {
	c.log.Debugf("Starting updating function '%s'", name)
	defer c.log.Debugf("Update complete for function '%s'", name)

	// TODO update function using deploy

	c.log.Debugf("Successfully updated function: '%s'", name)
	return nil
}

// Remove removes a stack and unregister any resources created.
func (c *CLIManager) Remove(name string) error {
	c.log.Debugf("Removing function: %s", name)
	defer c.log.Debugf("Removal of function '%s' complete", name)

	// TODO remove function

	return nil
}

// NewCLI returns the interface to manage function on Amazon lambda.
func NewCLI(
	log *logp.Logger,
	cfg *common.Config,
	provider provider.Provider,
) (provider.CLIManager, error) {
	config := &Config{}
	if err := cfg.Unpack(config); err != nil {
		return nil, err
	}

	builder, err := provider.TemplateBuilder()
	if err != nil {
		return nil, err
	}

	templateBuilder, ok := builder.(*restAPITemplateBuilder)
	if !ok {
		return nil, fmt.Errorf("not restAPITemplateBuilder")
	}

	return &CLIManager{
		config:          config,
		log:             logp.NewLogger("gcp"),
		templateBuilder: templateBuilder,
	}, nil
}
