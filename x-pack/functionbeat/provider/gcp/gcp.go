// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package gcp

import (
	"context"

	"cloud.google.com/go/pubsub"

	_ "github.com/elastic/beats/x-pack/functionbeat/provider/gcp/include"
)

func RunPubSub(ctx context.Context, m pubsub.Message) {
}
