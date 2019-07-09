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

type opDeleteFunction struct {
	log  *logp.Logger
	name string
}

func newOpDeleteFunction(log *logp.Logger, name string) *opDeleteFunction {
	return &opDeleteFunction{log: log}
}

func (o *opDeleteFunction) Execute(_ executor.Context) error {
	functionURL := googleAPIsURL + name
	req, err := http.Request("DELETE", functionURL, nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)

	fmt.Println(resp)

	return err
}
