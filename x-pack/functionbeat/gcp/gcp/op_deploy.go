package gcp

import "github.com/elastic/beats/libbeat/logp"

type opDeployFunction struct {
	log *logp.Logger
}

func newOpDeployFunction() *opDeployFunction {
	return &opDeployFunction{}
}

func (o *opDeployFunction) Execute(ctx executionContext) error {

}
