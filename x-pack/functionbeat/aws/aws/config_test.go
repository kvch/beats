// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucket(t *testing.T) {
	t.Run("valid bucket name", func(t *testing.T) {
		b := bucket("")
		err := b.Unpack("hello-bucket")
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, bucket("hello-bucket"), b)
	})

	t.Run("too long", func(t *testing.T) {
		b := bucket("")
		err := b.Unpack("hello-bucket-abc12345566789012345678901234567890123456789012345678901234567890")
		assert.Error(t, err)
	})

	t.Run("too short", func(t *testing.T) {
		b := bucket("")
		err := b.Unpack("he")
		assert.Error(t, err)
	})
}
