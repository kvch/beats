// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package include

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("heartbeat/fields.yml", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}
