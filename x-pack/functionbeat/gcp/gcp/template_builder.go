// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package gcp

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/x-pack/functionbeat/function/provider"
)

// NewTemplateBuilder returns the requested template builder
func NewTemplateBuilder(log *logp.Logger, cfg *common.Config, p provider.Provider) (provider.TemplateBuilder, error) {
	return newRestAPITemplateBuilder(log, cfg, p)
}

// restAPITemplateBuilder builds request object when deploying Functionbeat using
// the command deploy.
type restAPITemplateBuilder struct {
}

// newRestAPITemplateBuilder
func newRestAPITemplateBuilder(log *logp.Logger, cfg *common.Config, p provider.Provider) (provider.TemplateBuilder, error) {
	return &restAPITemplateBuilder{}, nil
}

func (r *restAPITemplateBuilder) getRequestBody() common.MapStr {
	return common.MapStr{}
}

// RawTemplate returns the JSON to POST to the endpoint.
func (r *restAPITemplateBuilder) RawTemplate(name string) (string, error) {
	return "", nil
}

// deploymentManaegerTemplateBuilder builds a YAML configuration for users
// to deploy the exported configuration using Google Deployment Manager.
type deploymentManaegerTemplateBuilder struct {
}

// newDeploymentManagerTemplateBuilder
func newDeploymentManagerTemplateBuilder(log *logp.Logger, cfg *common.Config, p provider.Provider) (provider.TemplateBuilder, error) {
	return &deploymentManaegerTemplateBuilder{}, nil
}

// RawTemplate returns YAML representation of the function to be deployed.
func (d *deploymentManaegerTemplateBuilder) RawTemplate(name string) (string, error) {
	return "", nil
}
