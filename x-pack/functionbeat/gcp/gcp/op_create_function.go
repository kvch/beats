// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package gcp

import (
	"fmt"
	"net/http"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/x-pack/functionbeat/function/executor"
)

type opCreateFunction struct {
	log             *logp.Logger
	templateBuilder *restAPITemplateBuilder
}

func newOpCreateFunction(log *logp.Logger, templateBuilder *restAPITemplateBuilder) *opCreateFunction {
	return &opCreateFunction{log: log, templateBuilder: templateBuilder}
}

func (o *opCreateFunction) Execute(_ executor.Context) error {
	deployURL := googleAPIsURL + c.location + "/functions"
	body := o.templateBuilder.requestBody()
	resp, err := http.Post(deployURL, "application/json", body)

	fmt.Println(resp)
	return err
}

func (o *opCreateFunction) Rollback(_ executor.Context) error {

}
