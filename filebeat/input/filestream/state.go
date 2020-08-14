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

package filestream

import (
	"fmt"

	"github.com/elastic/beats/v7/libbeat/common/file"
)

type state struct {
	Source         string       `json:"source" struct:"source"`
	Offset         int64        `json:"offset" struct:"offset"`
	FileStateOS    file.StateOS `json:"file_state_os" struct:"file_state_os"`
	IdentifierName string       `json:"identifier_name" struct:"identifier_name"`
}

func (s *state) String() string {
	return fmt.Sprintf(
		"{Source: %v, Offset: %v, FileStateOS: %v, IdentifierName: %v}",
		s.Source,
		s.Offset,
		s.FileStateOS,
		s.IdentifierName,
	)
}
