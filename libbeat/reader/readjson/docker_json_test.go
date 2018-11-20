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

package readjson

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/reader"
)

func TestDockerJSON(t *testing.T) {
	tests := []struct {
		name            string
		input           [][]byte
		stream          string
		partial         bool
		forceCRI        bool
		criflags        bool
		keepOriginal    bool
		expectedError   bool
		expectedMessage reader.Message
	}{
		{
			name:   "Common log message",
			input:  [][]byte{[]byte(`{"log":"1:M 09 Nov 13:27:36.276 # User requested shutdown...\n","stream":"stdout","time":"2017-11-09T13:27:36.277747246Z"}`)},
			stream: "all",
			expectedMessage: reader.Message{
				Content: []byte("1:M 09 Nov 13:27:36.276 # User requested shutdown...\n"),
				Fields:  common.MapStr{"stream": "stdout"},
				Ts:      time.Date(2017, 11, 9, 13, 27, 36, 277747246, time.UTC),
				Bytes:   122,
			},
		},
		{
			name:          "Wrong JSON",
			input:         [][]byte{[]byte(`this is not JSON`)},
			stream:        "all",
			expectedError: true,
		},
		{
			name:          "Wrong CRI",
			input:         [][]byte{[]byte(`2017-09-12T22:32:21.212861448Z stdout`)},
			stream:        "all",
			expectedError: true,
		},
		{
			name:          "Wrong CRI",
			input:         [][]byte{[]byte(`{this is not JSON nor CRI`)},
			stream:        "all",
			expectedError: true,
		},
		{
			name:          "Missing time",
			input:         [][]byte{[]byte(`{"log":"1:M 09 Nov 13:27:36.276 # User requested shutdown...\n","stream":"stdout"}`)},
			stream:        "all",
			expectedError: true,
		},
		{
			name:   "CRI log no tags",
			input:  [][]byte{[]byte(`2017-09-12T22:32:21.212861448Z stdout 2017-09-12 22:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache`)},
			stream: "all",
			expectedMessage: reader.Message{
				Content:  []byte("2017-09-12 22:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache"),
				Original: []byte("2017-09-12T22:32:21.212861448Z stdout 2017-09-12 22:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache"),
				Fields:   common.MapStr{"stream": "stdout"},
				Ts:       time.Date(2017, 9, 12, 22, 32, 21, 212861448, time.UTC),
				Bytes:    115,
			},
			criflags:     false,
			keepOriginal: true,
		},
		{
			name:   "CRI log",
			input:  [][]byte{[]byte(`2017-09-12T22:32:21.212861448Z stdout F 2017-09-12 22:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache`)},
			stream: "all",
			expectedMessage: reader.Message{
				Content: []byte("2017-09-12 22:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache"),
				Fields:  common.MapStr{"stream": "stdout"},
				Ts:      time.Date(2017, 9, 12, 22, 32, 21, 212861448, time.UTC),
				Bytes:   117,
			},
			criflags: true,
		},
		{
			name: "Filtering stream",
			input: [][]byte{
				[]byte(`{"log":"filtered\n","stream":"stdout","time":"2017-11-09T13:27:36.277747246Z"}`),
				[]byte(`{"log":"unfiltered\n","stream":"stderr","time":"2017-11-09T13:27:36.277747246Z"}`),
				[]byte(`{"log":"unfiltered\n","stream":"stdout","time":"2017-11-09T13:27:36.277747246Z"}`),
			},
			stream: "stderr",
			expectedMessage: reader.Message{
				Content: []byte("unfiltered\n"),
				Fields:  common.MapStr{"stream": "stderr"},
				Ts:      time.Date(2017, 11, 9, 13, 27, 36, 277747246, time.UTC),
				Bytes:   80,
			},
		},
		{
			name: "Filtering stream",
			input: [][]byte{
				[]byte(`2017-10-12T13:32:21.232861448Z stdout F 2017-10-12 13:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache`),
				[]byte(`2017-11-12T23:32:21.212771448Z stderr F 2017-11-12 23:32:21.212 [ERROR][77] table.go 111: error`),
				[]byte(`2017-12-12T10:32:21.212864448Z stdout F 2017-12-12 10:32:21.212 [WARN][88] table.go 222: Warn`),
			},
			stream: "stderr",
			expectedMessage: reader.Message{
				Content: []byte("2017-11-12 23:32:21.212 [ERROR][77] table.go 111: error"),
				Fields:  common.MapStr{"stream": "stderr"},
				Ts:      time.Date(2017, 11, 12, 23, 32, 21, 212771448, time.UTC),
				Bytes:   95,
			},
			criflags: true,
		},
		{
			name: "Split lines",
			input: [][]byte{
				[]byte(`{"log":"1:M 09 Nov 13:27:36.276 # User requested ","stream":"stdout","time":"2017-11-09T13:27:36.277747246Z"}`),
				[]byte(`{"log":"shutdown...\n","stream":"stdout","time":"2017-11-09T13:27:36.277747246Z"}`),
			},
			stream:  "stdout",
			partial: true,
			expectedMessage: reader.Message{
				Content: []byte("1:M 09 Nov 13:27:36.276 # User requested shutdown...\n"),
				Fields:  common.MapStr{"stream": "stdout"},
				Ts:      time.Date(2017, 11, 9, 13, 27, 36, 277747246, time.UTC),
				Bytes:   190,
			},
		},
		{
			name: "CRI Split lines",
			input: [][]byte{
				[]byte(`2017-10-12T13:32:21.232861448Z stdout P 2017-10-12 13:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache`),
				[]byte(`2017-11-12T23:32:21.212771448Z stdout F  error`),
			},
			stream:  "stdout",
			partial: true,
			expectedMessage: reader.Message{
				Content:  []byte("2017-10-12 13:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache error"),
				Original: []byte("2017-10-12T13:32:21.232861448Z stdout P 2017-10-12 13:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache\n2017-11-12T23:32:21.212771448Z stdout F  error"),
				Fields:   common.MapStr{"stream": "stdout"},
				Ts:       time.Date(2017, 10, 12, 13, 32, 21, 232861448, time.UTC),
				Bytes:    163,
			},
			criflags:     true,
			keepOriginal: true,
		},
		{
			name: "Split lines and remove \\n",
			input: [][]byte{
				[]byte("2017-10-12T13:32:21.232861448Z stdout P 2017-10-12 13:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache\n"),
				[]byte("2017-11-12T23:32:21.212771448Z stdout F  error"),
			},
			stream: "stdout",
			expectedMessage: reader.Message{
				Content: []byte("2017-10-12 13:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache error"),
				Fields:  common.MapStr{"stream": "stdout"},
				Ts:      time.Date(2017, 10, 12, 13, 32, 21, 232861448, time.UTC),
				Bytes:   164,
			},
			partial:  true,
			criflags: true,
		},
		{
			name: "Split lines with partial disabled",
			input: [][]byte{
				[]byte(`{"log":"1:M 09 Nov 13:27:36.276 # User requested ","stream":"stdout","time":"2017-11-09T13:27:36.277747246Z"}`),
				[]byte(`{"log":"shutdown...\n","stream":"stdout","time":"2017-11-09T13:27:36.277747246Z"}`),
			},
			stream:  "stdout",
			partial: false,
			expectedMessage: reader.Message{
				Content: []byte("1:M 09 Nov 13:27:36.276 # User requested "),
				Fields:  common.MapStr{"stream": "stdout"},
				Ts:      time.Date(2017, 11, 9, 13, 27, 36, 277747246, time.UTC),
				Bytes:   109,
			},
		},
		{
			name:          "Force CRI with JSON logs",
			input:         [][]byte{[]byte(`{"log":"1:M 09 Nov 13:27:36.276 # User requested shutdown...\n","stream":"stdout"}`)},
			stream:        "all",
			forceCRI:      true,
			expectedError: true,
		},
		{
			name:   "Force CRI log no tags",
			input:  [][]byte{[]byte(`2017-09-12T22:32:21.212861448Z stdout 2017-09-12 22:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache`)},
			stream: "all",
			expectedMessage: reader.Message{
				Content: []byte("2017-09-12 22:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache"),
				Fields:  common.MapStr{"stream": "stdout"},
				Ts:      time.Date(2017, 9, 12, 22, 32, 21, 212861448, time.UTC),
				Bytes:   115,
			},
			forceCRI: true,
			criflags: false,
		},
		{
			name:   "Force CRI log with flags",
			input:  [][]byte{[]byte(`2017-09-12T22:32:21.212861448Z stdout F 2017-09-12 22:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache`)},
			stream: "all",
			expectedMessage: reader.Message{
				Content: []byte("2017-09-12 22:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache"),
				Fields:  common.MapStr{"stream": "stdout"},
				Ts:      time.Date(2017, 9, 12, 22, 32, 21, 212861448, time.UTC),
				Bytes:   117,
			},
			forceCRI: true,
			criflags: true,
		},
		{
			name: "Force CRI split lines",
			input: [][]byte{
				[]byte(`2017-10-12T13:32:21.232861448Z stdout P 2017-10-12 13:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache`),
				[]byte(`2017-11-12T23:32:21.212771448Z stdout F  error`),
			},
			stream:  "stdout",
			partial: true,
			expectedMessage: reader.Message{
				Content: []byte("2017-10-12 13:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache error"),
				Fields:  common.MapStr{"stream": "stdout"},
				Ts:      time.Date(2017, 10, 12, 13, 32, 21, 232861448, time.UTC),
				Bytes:   163,
			},
			forceCRI: true,
			criflags: true,
		},
		{
			name: "Force CRI split lines and remove \\n",
			input: [][]byte{
				[]byte("2017-10-12T13:32:21.232861448Z stdout P 2017-10-12 13:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache\n"),
				[]byte("2017-11-12T23:32:21.212771448Z stdout F  error"),
			},
			stream: "stdout",
			expectedMessage: reader.Message{
				Content: []byte("2017-10-12 13:32:21.212 [INFO][88] table.go 710: Invalidating dataplane cache error"),
				Fields:  common.MapStr{"stream": "stdout"},
				Ts:      time.Date(2017, 10, 12, 13, 32, 21, 232861448, time.UTC),
				Bytes:   164,
			},
			partial:  true,
			forceCRI: true,
			criflags: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := &mockReader{messages: test.input}
			json := New(r, test.stream, test.partial, test.forceCRI, test.criflags, test.keepOriginal)
			message, err := json.Next()

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err == nil {
				assert.EqualValues(t, test.expectedMessage, message)
			}
		})
	}
}

type mockReader struct {
	messages [][]byte
}

func (m *mockReader) Next() (reader.Message, error) {
	message := m.messages[0]
	m.messages = m.messages[1:]
	return reader.Message{
		Content: message,
		Bytes:   len(message),
	}, nil
}
