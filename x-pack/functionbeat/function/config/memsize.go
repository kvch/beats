// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package config

import (
	"fmt"
	"unicode"

	humanize "github.com/dustin/go-humanize"

	"github.com/elastic/beats/libbeat/common/cfgwarn"
)

// MemSizeFactor64 implements a human understandable format for bytes but also make sure that all
// values used are a factory of 64.
type MemSizeFactor64 int

// Unpack converts a size defined from a human readable format into bytes and ensure that the value
// is a factoru of 64.
func (m *MemSizeFactor64) Unpack(v string) error {
	sz, err := humanize.ParseBytes(v)
	if isRawBytes(v) {
		cfgwarn.Deprecate("7.0", "size now requires a unit (KiB, MiB, etc...), current value: %s.", v)
	}
	if err != nil {
		return err
	}

	if sz%64 != 0 {
		return fmt.Errorf("number is not a factor of 64, %d bytes (user value: %s)", sz, v)
	}

	*m = MemSizeFactor64(sz)
	return nil
}

// Megabytes return the value in megatebytes.
func (m *MemSizeFactor64) Megabytes() int {
	return int(*m) / 1024 / 1024
}

func isRawBytes(v string) bool {
	for _, c := range v {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
