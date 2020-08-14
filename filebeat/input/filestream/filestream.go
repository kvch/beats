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
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/elastic/beats/v7/filebeat/harvester"
	input "github.com/elastic/beats/v7/filebeat/input/v2"
	"github.com/elastic/beats/v7/libbeat/common/backoff"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/go-concert/ctxtool"
	"github.com/elastic/go-concert/unison"
)

var (
	ErrFileTruncate = errors.New("detected file being truncated")
	ErrClosed       = errors.New("reader closed")
)

// Log contains all log related data
type Log struct {
	fs            harvester.Source
	log           *logp.Logger
	ctx           context.Context
	cancelReading context.CancelFunc

	closeInactive time.Duration
	closeTimeout  time.Duration
	closeEOF      bool

	offset       int64
	lastTimeRead time.Time
	backoff      backoff.Backoff
	tg           unison.TaskGroup
}

// NewLog creates a new log instance to read log sources
func newFileReader(
	log *logp.Logger,
	canceler input.Canceler,
	fs harvester.Source,
	config readerConfig,
) (*Log, error) {
	var offset int64
	if seeker, ok := fs.(io.Seeker); ok {
		var err error
		offset, err = seeker.Seek(0, os.SEEK_CUR)
		if err != nil {
			return nil, err
		}
	}

	ctx, cancelReading := context.WithCancel(ctxtool.FromCanceller(canceler))

	l := &Log{
		fs:            fs,
		log:           log,
		ctx:           ctx,
		cancelReading: cancelReading,
		offset:        offset,
		lastTimeRead:  time.Now(),
		backoff:       backoff.NewExpBackoff(canceler.Done(), config.Backoff.Init, config.Backoff.Max),
	}

	l.startFileMonitoringIfNeeded()

	return l, nil
}

// Read reads from the reader and updates the offset
// The total number of bytes read is returned.
func (f *Log) Read(buf []byte) (int, error) {
	totalN := 0

	for f.ctx.Err() == nil {
		n, err := f.fs.Read(buf)
		if n > 0 {
			f.offset += int64(n)
			f.lastTimeRead = time.Now()
		}
		totalN += n

		// Read from source completed without error
		// Either end reached or buffer full
		if err == nil {
			// reset backoff for next read
			f.backoff.Reset()
			return totalN, nil
		}

		// Move buffer forward for next read
		buf = buf[n:]

		// Checks if an error happened or buffer is full
		// If buffer is full, cannot continue reading.
		// Can happen if n == bufferSize + io.EOF error
		err = f.errorChecks(err)
		if err != nil || len(buf) == 0 {
			return totalN, err
		}

		f.log.Debugf("End of file reached: %s; Backoff now.", f.fs.Name())
		f.backoff.Wait()
	}

	return 0, ErrClosed
}

func (f *Log) startFileMonitoringIfNeeded() {
	if f.closeInactive == 0 && f.closeTimeout == 0 {
		return
	}

	f.tg = unison.TaskGroup{}

	if f.closeInactive > 0 {
		f.tg.Go(func(ctx unison.Canceler) error {
			f.closeIfTimeout(ctx)
			return nil
		})
	}

	if f.closeTimeout > 0 {
		f.tg.Go(func(ctx unison.Canceler) error {
			f.closeIfInactive(ctx)
			return nil
		})
	}
}

func (f *Log) closeIfTimeout(ctx unison.Canceler) {
	timer := time.NewTimer(f.closeTimeout)
	for ctx.Err() == nil {
		select {
		case <-timer.C:
			f.cancelReading()
			return
		}
	}
	f.log.Debug("Monitoring if timeout has been reached has ended.")
}

func (f *Log) closeIfInactive(ctx unison.Canceler) {
	for ctx.Err() == nil {
		age := time.Since(f.lastTimeRead)
		if age > f.closeInactive {
			f.cancelReading()
			return
		}
	}
	f.log.Debug("Monitoring if file is inactive has ended.")
}

// errorChecks determines the cause for EOF errors, and how the EOF event should be handled
// based on the config options.
func (f *Log) errorChecks(err error) error {
	if err != io.EOF {
		f.log.Error("Unexpected state reading from %s; error: %s", f.fs.Name(), err)
		return err
	}

	return f.handleEOF()
}

func (f *Log) handleEOF() error {
	err := io.EOF

	if f.closeEOF {
		return err
	}

	// Refetch fileinfo to check if the file was truncated.
	// Errors if the file was removed/rotated after reading and before
	// calling the stat function
	info, statErr := f.fs.Stat()
	if statErr != nil {
		f.log.Error("Unexpected error reading from %s; error: %s", f.fs.Name(), statErr)
		return statErr
	}

	// check if file was truncated
	if info.Size() < f.offset {
		f.log.Debugf("File was truncated as offset (%d) > size (%d): %s", f.offset, info.Size(), f.fs.Name())
		return ErrFileTruncate
	}

	return nil
}

// Close
func (f *Log) Close() error {
	// Note: File reader is not closed here because that leads to race conditions
	f.cancelReading()
	return f.tg.Stop()
}
