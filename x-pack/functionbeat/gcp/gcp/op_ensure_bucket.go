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
	config Config
	name   string
}

func newOpEnsureBucket(log *logp.Logger, cfg Config, name string) *opEnsureBucket {
	return &opEnsureBucket{log: log, config: cfg, name: name}
}

func (o *opEnsureBucket) Execute(_ executor.Context) error {
	o.log.Debugf("Verifying presence of Cloud Storage bucket: %s", o.name)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucket := client.Bucket(name)
	err = bucket.Attrs(ctx)
	if gerr, ok := err.(storage.Error); ok {
		if gerr == storage.ErrBucketNotExist {
			err = bucket.Create(ctx, o.config.ProjectID, nil)
			if err != nil {
				return fmt.Errorf("cannot create bucket for function: %+v", err)
			}
			o.log.Debugf("Cloud Storage bucket created with name '%s'", o.name)
		}
	}

	return fmt.Errorf("cannot get info on bucket and does exist +%v", err)
}
