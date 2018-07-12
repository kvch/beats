// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package encode

import (
	"io"
	"time"

	"github.com/elastic/beats/filebeat/reader"
	"github.com/elastic/beats/filebeat/reader/encode/encoding"
	"github.com/elastic/beats/filebeat/reader/line"
	"github.com/elastic/beats/libbeat/common"
)

// Reader produces lines by reading lines from an io.Reader
// through a decoder converting the reader it's encoding to utf-8.
type Reader struct {
	reader *line.Reader
}

// New creates a new Encode reader from input reader by applying
// the given codec.
func New(
	r io.Reader,
	codec encoding.Encoding,
	bufferSize int,
) Reader {
	eReader := line.New(r, codec, []byte("\n"), bufferSize)
	return Reader{eReader}
}

// Next reads the next line from it's initial io.Reader
// This converts a io.Reader to a reader.reader
func (r Reader) Next() (reader.Message, error) {
	c, sz, err := r.reader.Next()
	// Creating message object
	return reader.Message{
		Ts:      time.Now(),
		Content: c,
		Bytes:   sz,
	}, err
}

func (r Reader) GetState() common.MapStr {
	return r.reader.GetState()
}
