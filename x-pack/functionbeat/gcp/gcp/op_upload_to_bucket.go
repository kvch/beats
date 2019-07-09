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

type opUploadToBucket struct {
	log    *logp.Logger
	config *Config
	raw    []byte
}

func newOpUploadToBucket(log *logp.Logger, config *Config, raw []byte) *opUploadToBucket {
	return &opUploadToBucket{
		log:    log,
		config: config,
		raw:    raw,
	}
}

func (o *opUploadToBucket) Execute(_ executor.Context) error {
	o.log.Debugf("Uploading file 'functionbeat-gcp' to bucket '%s' with size %d bytes", o.config.FunctionStorage, len(o.raw))

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("could not create storage client: %+v", err)
	}
	w := client.Bucket(o.config.FunctionStorage).Object("functionbeat-gcp").NewWriter(ctx)
	w.ContentType = "text/plain"
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}} // TODO check permissions
	_, err = w.Write(o.raw)
	if err != nil {
		return fmt.Errorf("error while writing function: %+v", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("error while closing writer: %+v", err)
	}

	o.log.Debug("Upload successful", w.Attrs())
	return nil
}

// TODO
func (o *opUploadToBucket) Rollback(ctx executor.Context) error {
	return nil
}
