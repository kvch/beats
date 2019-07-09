// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/x-pack/functionbeat/function/executor"
)

type opEnsureBucket struct {
	log    *logp.Logger
	config *Config
}

func newOpEnsureBucket(log *logp.Logger, cfg *Config) *opEnsureBucket {
	return &opEnsureBucket{log: log, config: cfg}
}

func (o *opEnsureBucket) Execute(_ executor.Context) error {
	o.log.Debugf("Verifying presence of Cloud Storage bucket: %s", o.config.FunctionStorage)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucket := client.Bucket(o.config.FunctionStorage)
	attrs, err := bucket.Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		berr := bucket.Create(ctx, o.config.ProjectID, nil)
		if berr != nil {
			return fmt.Errorf("cannot create bucket for function: %+v", berr)
		}
		o.log.Debugf("Cloud Storage bucket created with name '%s', attributes: %+v", o.config.FunctionStorage, attrs)
	}

	return fmt.Errorf("cannot get info on bucket and does exist +%v", err)
}
