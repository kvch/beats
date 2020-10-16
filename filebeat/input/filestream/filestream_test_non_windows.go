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

// +build !windows

package filestream

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/elastic/beats/v7/libbeat/logp"
)

// these tests are separated as one cannot delete/rename files
// while another process is working with it on Windows
func TestLogFileRenamed(t *testing.T) {
	f := createTestLogFile()
	defer f.Close()

	renamedFile := f.Name() + ".renamed"

	reader, err := newFileReader(
		logp.L(),
		context.TODO(),
		f,
		readerConfig{},
		closerConfig{
			OnStateChange: stateChangeCloserConfig{
				CheckInterval: 1 * time.Second,
				Renamed:       true,
			},
		},
	)
	if err != nil {
		t.Fatalf("error while creating logReader: %+v", err)
	}

	buf := make([]byte, 1024)
	_, err = reader.Read(buf)
	assert.Nil(t, err)

	err = os.Rename(f.Name(), renamedFile)
	if err != nil {
		t.Fatalf("error while renaming file: %+v", err)
	}

	err = readUntilError(reader)
	os.Remove(renamedFile)

	assert.Equal(t, ErrClosed, err)
}

func TestLogFileRemoved(t *testing.T) {
	f := createTestLogFile()
	defer f.Close()

	reader, err := newFileReader(
		logp.L(),
		context.TODO(),
		f,
		readerConfig{},
		closerConfig{
			OnStateChange: stateChangeCloserConfig{
				CheckInterval: 1 * time.Second,
				Removed:       true,
			},
		},
	)
	if err != nil {
		t.Fatalf("error while creating logReader: %+v", err)
	}

	buf := make([]byte, 1024)
	_, err = reader.Read(buf)
	assert.Nil(t, err)

	err = os.Remove(f.Name())
	if err != nil {
		t.Fatalf("error while remove file: %+v", err)
	}

	err = readUntilError(reader)

	assert.Equal(t, ErrClosed, err)
}
