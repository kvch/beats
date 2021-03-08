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

// +build integration

package filestream

import (
	"context"
	"os"
	"runtime"
	"testing"

	loginp "github.com/elastic/beats/v7/filebeat/input/filestream/internal/input-logfile"
)

// test_close_renamed from test_harvester.py
func TestFilestreamCloseRenamed(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("renaming files while Filebeat is running is not supported on Windows")
	}

	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths":                                []string{env.abspath(testlogName) + "*"},
		"prospector.scanner.check_interval":    "1ms",
		"close.on_state_change.check_interval": "1ms",
		"close.on_state_change.renamed":        "true",
	})

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	testlines := []byte("first log line\n")
	env.mustWriteLinesToFile(testlogName, testlines)

	// first event has made it successfully
	env.waitUntilEventCount(1)
	env.requireOffsetInRegistry(testlogName, len(testlines))

	testlogNameRotated := "test.log.rotated"
	env.mustRenameFile(testlogName, testlogNameRotated)

	newerTestlines := []byte("new first log line\nnew second log line\n")
	env.mustWriteLinesToFile(testlogName, newerTestlines)

	// new two events arrived
	env.waitUntilEventCount(3)

	cancelInput()
	env.waitUntilInputStops()

	env.requireOffsetInRegistry(testlogNameRotated, len(testlines))
	env.requireOffsetInRegistry(testlogName, len(newerTestlines))
}

// test_close_removed from test_harvester.py
func TestFilestreamCloseRemoved(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths":                                []string{env.abspath(testlogName) + "*"},
		"prospector.scanner.check_interval":    "24h",
		"close.on_state_change.check_interval": "1ms",
		"close.on_state_change.removed":        "true",
	})

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	testlines := []byte("first log line\n")
	env.mustWriteLinesToFile(testlogName, testlines)

	// first event has made it successfully
	env.waitUntilEventCount(1)
	// check registry
	env.requireOffsetInRegistry(testlogName, len(testlines))

	fi, err := os.Stat(env.abspath(testlogName))
	if err != nil {
		t.Fatalf("cannot stat file: %+v", err)
	}

	env.mustRemoveFile(testlogName)

	// the second log line will not be picked up as scan_interval is set to one day.
	env.mustWriteLinesToFile(testlogName, []byte("first line\nsecond log line\n"))

	// new two events arrived
	env.waitUntilEventCount(1)

	cancelInput()
	env.waitUntilInputStops()

	identifier, _ := newINodeDeviceIdentifier(nil)
	src := identifier.GetSource(loginp.FSEvent{Info: fi, Op: loginp.OpCreate, NewPath: env.abspath(testlogName)})
	env.requireOffsetInRegistryByID(src.Name(), len(testlines))
}

// test_close_eof from test_harvester.py
func TestFilestreamCloseEOF(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths":                             []string{env.abspath(testlogName)},
		"prospector.scanner.check_interval": "24h",
		"close.reader.on_eof":               "true",
	})

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	testlines := []byte("first log line\n")
	expectedOffset := len(testlines)
	env.mustWriteLinesToFile(testlogName, testlines)

	// first event has made it successfully
	env.waitUntilEventCount(1)
	env.requireOffsetInRegistry(testlogName, expectedOffset)

	// the second log line will not be picked up as scan_interval is set to one day.
	env.mustWriteLinesToFile(testlogName, []byte("first line\nsecond log line\n"))

	// only one event is read
	env.waitUntilEventCount(1)

	cancelInput()
	env.waitUntilInputStops()

	env.requireOffsetInRegistry(testlogName, expectedOffset)
}

// test_empty_lines from test_harvester.py
func TestFilestreamEmptyLine(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths":                             []string{env.abspath(testlogName)},
		"prospector.scanner.check_interval": "1ms",
	})

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	testlines := []byte("first log line\nnext is an empty line\n")
	env.mustWriteLinesToFile(testlogName, testlines)

	env.waitUntilEventCount(2)
	env.requireOffsetInRegistry(testlogName, len(testlines))

	moreTestlines := []byte("\nafter an empty line\n")
	env.mustAppendLinesToFile(testlogName, moreTestlines)

	env.waitUntilEventCount(3)

	cancelInput()
	env.waitUntilInputStops()

	env.requireOffsetInRegistry(testlogName, len(testlines)+len(moreTestlines))
}

// test_empty_lines_only from test_harvester.py
// This test differs from the original because in filestream
// input offset is no longer persisted when the line is empty.
func TestFilestreamEmptyLinesOnly(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths":                             []string{env.abspath(testlogName)},
		"prospector.scanner.check_interval": "1ms",
	})

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	testlines := []byte("\n\n\n")
	env.mustWriteLinesToFile(testlogName, testlines)

	cancelInput()
	env.waitUntilInputStops()

	env.requireNoEntryInRegistry(testlogName)
}

// test_exceed_buffer from test_harvester.py
func TestFilestreamExceedBuffer(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths":       []string{env.abspath(testlogName)},
		"buffer_size": 10,
	})

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	testline := []byte("a line longer than size allowed in buffer_size\n")
	expectedOffset := len(testline)
	env.mustWriteLinesToFile(testlogName, testline)

	// event arrives to the output in full
	env.waitUntilEventCount(1)
	env.requireEventsReceived([]string{string(testline[:len(testline)-1])})

	cancelInput()
	env.waitUntilInputStops()

	env.requireOffsetInRegistry(testlogName, expectedOffset)
}

// test_truncated_file_open from test_harvester.py
func TestFilestreamTruncatedFileOpen(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths":                             []string{env.abspath(testlogName)},
		"prospector.scanner.check_interval": "1ms",
	})

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	testlines := []byte("first line\nsecond line\nthird line\n")
	env.mustWriteLinesToFile(testlogName, testlines)

	env.waitUntilEventCount(3)
	env.requireOffsetInRegistry(testlogName, len(testlines))

	env.mustTruncateFile(testlogName, 0)

	truncatedTestLines := []byte("truncated first line\n")
	env.mustWriteLinesToFile(testlogName, truncatedTestLines)
	env.waitUntilEventCount(4)

	cancelInput()
	env.waitUntilInputStops()
	env.requireOffsetInRegistry(testlogName, len(truncatedTestLines))
}

// test_truncated_file_closed from test_harvester.py
func TestFilestreamTruncatedFileClosed(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths":                                []string{env.abspath(testlogName)},
		"prospector.scanner.check_interval":    "1ms",
		"close.on_state_change.check_interval": "1ms",
		"close.on_state_change.inactive":       "50ms",
	})

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	testlines := []byte("first line\nsecond line\nthird line\n")
	env.mustWriteLinesToFile(testlogName, testlines)

	env.waitUntilEventCount(3)
	env.requireOffsetInRegistry(testlogName, len(testlines))

	env.waitUntilHarvesterIsDone()

	env.mustTruncateFile(testlogName, 0)

	truncatedTestLines := []byte("truncated first line\n")
	env.mustWriteLinesToFile(testlogName, truncatedTestLines)
	env.waitUntilEventCount(4)

	cancelInput()
	env.waitUntilInputStops()
	env.requireOffsetInRegistry(testlogName, len(truncatedTestLines))
}

// test_truncated_file_closed from test_harvester.py
func TestFilestreamCloseTimeout(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths":                                []string{env.abspath(testlogName)},
		"prospector.scanner.check_interval":    "24h",
		"close.on_state_change.check_interval": "100ms",
		"close.reader.after_interval":          "500ms",
	})

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	testlines := []byte("first line\n")
	env.mustWriteLinesToFile(testlogName, testlines)

	env.waitUntilEventCount(1)
	env.requireOffsetInRegistry(testlogName, len(testlines))
	env.waitUntilHarvesterIsDone()

	env.mustWriteLinesToFile(testlogName, []byte("first line\nsecond log line\n"))

	env.waitUntilEventCount(1)

	cancelInput()
	env.waitUntilInputStops()

	env.requireOffsetInRegistry(testlogName, len(testlines))
}


// test_symlinks_enabled from test_harvester.py
func TestFilestreamSymlinksEnabled(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	symlinkName := "test.log.symlink"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths": []string{
			env.abspath(testlogName),
			env.abspath(symlinkName),
		},
		"prospector.scanner.symlinks": "true",
	})

	testlines := []byte("first line\n")
	env.mustWriteLinesToFile(testlogName, testlines)

	env.mustSymlink(testlogName, symlinkName)

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	env.waitUntilEventCount(1)

	cancelInput()
	env.waitUntilInputStops()

	env.requireOffsetInRegistry(testlogName, len(testlines))
	env.requireOffsetInRegistry(symlinkName, len(testlines))
}

// test_symlink_rotated from test_harvester.py
func TestFilestreamSymlinkRotated(t *testing.T) {
	env := newInputTestingEnvironment(t)

	firstTestlogName := "test1.log"
	secondTestlogName := "test2.log"
	symlinkName := "test.log"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths": []string{
			env.abspath(firstTestlogName),
			env.abspath(secondTestlogName),
			env.abspath(symlinkName),
		},
		"prospector.scanner.check_interval": "1ms",
		"prospector.scanner.symlinks":       "true",
		"close.on_state_change.removed":     "false",
		"clean_removed":                     "false",
	})

	commonLine := "first line in file "
	for i, path := range []string{firstTestlogName, secondTestlogName} {
		env.mustWriteLinesToFile(path, []byte(commonLine+strconv.Itoa(i)+"\n"))
	}

	env.mustSymlink(firstTestlogName, symlinkName)

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	env.waitUntilEventCount(2)

	expectedOffset := len(commonLine) + 2
	env.requireOffsetInRegistry(firstTestlogName, expectedOffset)
	env.requireOffsetInRegistry(secondTestlogName, expectedOffset)

	// rotate symlink
	env.mustRemoveFile(symlinkName)
	env.mustSymlink(secondTestlogName, symlinkName)

	moreLines := "second line in file 2\nthird line in file 2\n"
	env.mustAppendLinesToFile(secondTestlogName, []byte(moreLines))

	env.waitUntilEventCount(4)
	env.requireOffsetInRegistry(firstTestlogName, expectedOffset)
	env.requireOffsetInRegistry(secondTestlogName, expectedOffset+len(moreLines))
	env.requireOffsetInRegistry(symlinkName, expectedOffset+len(moreLines))

	cancelInput()
	env.waitUntilInputStops()

	env.requireRegistryEntryCount(2)
}

// test_symlink_removed from test_harvester.py
func TestFilestreamSymlinkRemoved(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	symlinkName := "test.log.symlink"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths": []string{
			env.abspath(testlogName),
			env.abspath(symlinkName),
		},
		"prospector.scanner.check_interval": "1ms",
		"prospector.scanner.symlinks":       "true",
		"close.on_state_change.removed":     "false",
		"clean_removed":                     "false",
	})

	line := []byte("first line\n")
	env.mustWriteLinesToFile(testlogName, line)

	env.mustSymlink(testlogName, symlinkName)

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	env.waitUntilEventCount(1)

	env.requireOffsetInRegistry(testlogName, len(line))

	// remove symlink
	env.mustRemoveFile(symlinkName)

	env.mustAppendLinesToFile(testlogName, line)

	env.waitUntilEventCount(2)
	env.requireOffsetInRegistry(testlogName, 2*len(line))

	cancelInput()
	env.waitUntilInputStops()

	env.requireRegistryEntryCount(1)
}

// test_truncate from test_harvester.py
func TestFilestreamTruncate(t *testing.T) {
	env := newInputTestingEnvironment(t)

	testlogName := "test.log"
	symlinkName := "test.log.symlink"
	inp := env.mustCreateInput(map[string]interface{}{
		"paths": []string{
			env.abspath(testlogName),
			env.abspath(symlinkName),
		},
		"prospector.scanner.check_interval": "1ms",
		"prospector.scanner.symlinks":       "true",
	})

	lines := []byte("first line\nsecond line\nthird line\n")
	env.mustWriteLinesToFile(testlogName, lines)

	env.mustSymlink(testlogName, symlinkName)

	ctx, cancelInput := context.WithCancel(context.Background())
	env.startInput(ctx, inp)

	env.waitUntilEventCount(3)

	env.requireOffsetInRegistry(testlogName, len(lines))

	// remove symlink
	env.mustRemoveFile(symlinkName)
	env.mustTruncateFile(testlogName, 0)

	moreLines := []byte("forth line\nfifth line\n")
	env.mustWriteLinesToFile(testlogName, moreLines)

	env.waitUntilEventCount(5)
	env.requireOffsetInRegistry(testlogName, len(moreLines))

	cancelInput()
	env.waitUntilInputStops()

	env.requireRegistryEntryCount(1)
}
