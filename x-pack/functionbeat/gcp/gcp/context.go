// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package gcp

import (
	"errors"
	"time"
)

var (
	errMissingStackID = errors.New("missing stack id")
)

type gcpContext struct {
	StartedAt time.Time
}

func newContext() *gcpContext {
	return &gcpContext{StartedAt: time.Now()}
}
