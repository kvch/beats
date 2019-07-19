// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"fmt"
	"regexp"
	"unicode"

	humanize "github.com/dustin/go-humanize"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/cfgwarn"
)

// ConfigOverrides overrides the defaults provided by libbeat.
var (
	functionPattern = "^[A-Za-z][A-Za-z0-9\\-]{0,139}$"
	functionRE      = regexp.MustCompile(functionPattern)
	ConfigOverrides = common.MustNewConfigFrom(map[string]interface{}{
		"path.data":              "/tmp",
		"path.logs":              "/tmp/logs",
		"logging.to_stderr":      true,
		"logging.to_files":       false,
		"setup.template.enabled": true,
		"queue.mem": map[string]interface{}{
			"events":           "${output.elasticsearch.bulk_max_size}",
			"flush.min_events": 10,
			"flush.timeout":    "0.01s",
		},
		"output.elasticsearch.bulk_max_size": 50,
	})
)

// Config default configuration for Functionbeat.
type Config struct {
	Provider *common.Config `config:"provider" validate:"required"`
}

// ProviderConfig is a generic configured used by providers.
type ProviderConfig struct {
	Functions []*common.Config `config:"functions"`
}

// FunctionConfig minimal configuration from each function.
type FunctionConfig struct {
	Type    string       `config:"type"`
	Name    functionName `config:"name"`
	Enabled bool         `config:"enabled"`
}

// DefaultConfig is the default configuration for Functionbeat.
var DefaultConfig = Config{}

// DefaultFunctionConfig is the default configuration for new function.
var DefaultFunctionConfig = FunctionConfig{
	Enabled: true,
}

type functionName string

func (f *functionName) Unpack(s string) error {
	if !functionRE.MatchString(s) {
		return fmt.Errorf(
			"invalid name: '%s', name must match [a-zA-Z0-9-] and be at most 140 characters",
			s,
		)
	}
	*f = functionName(s)
	return nil
}

func (f *functionName) String() string {
	return string(*f)
}

// Validate enforces that function names are unique.
func (p *ProviderConfig) Validate() error {
	names := make(map[functionName]bool)
	for _, rawfn := range p.Functions {
		fc := FunctionConfig{}
		rawfn.Unpack(&fc)

		if !fc.Enabled {
			return nil
		}

		if _, found := names[fc.Name]; found {
			return fmt.Errorf("function name '%s' already exist, name must be unique", fc.Name)
		}

		names[fc.Name] = true
	}
	return nil
}

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
