// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package gcp

import (
	"github.com/elastic/beats/libbeat/feature"
	"github.com/elastic/beats/x-pack/functionbeat/function/provider"
)

// Bundle exposes the trigger supported by the GCP provider.
var Bundle = provider.MustCreate(
	"gcp",
	provider.NewDefaultProvider("gcp", NewCLI, NewTemplateBuilder),
	feature.NewDetails("Google Cloud Platform", "listen to events from Google Cloud Platform", feature.Stable),
).MustAddFunction("pubsub",
	NewPubSub,
	feature.NewDetails(
		"Google Pub/Sub trigger",
		"receive events from Google Pub/Sub.",
		feature.Stable,
	),
).MustAddFunction("cloudstorage",
	NewCloudStorage,
	feature.NewDetails(
		"Google Cloud Storage trigger",
		"receive events from Google Cloud Storage.",
		feature.Stable,
	),
).Bundle()
