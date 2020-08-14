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
	"os"
	"path/filepath"
	"time"

	"github.com/elastic/beats/v7/filebeat/input/file"
	loginp "github.com/elastic/beats/v7/filebeat/input/filestream/internal/input-logfile"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/common/match"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/go-concert/unison"
)

const (
	recursiveGlobDepth = 8

	scannerName = "scanner"

	watcherDebugKey = "file_watcher"
	scannerDebugKey = "file_scanner"

	msgGlobFailed         = "glob(%s) failed: %v"
	msgExcludedFile       = "Exclude file: %s"
	msgLstatFileFailed    = "lstat(%s) failed: %s"
	msgDirSkipped         = "Skipping directory: %s"
	msgSymlinkSkipped     = "File %s skipped as it is a symlink"
	msgStatFileFailed     = "stat(%s) failed: %s"
	msgSymlinkSameSkipped = "Same file found as symlink and original. Skipping file: %s (as it same as %s)"
)

var (
	watcherFactories = map[string]watcherFactory{
		scannerName: newScannerWatcher,
	}
)

type watcherFactory func([]string, *common.Config) (loginp.FSWatcher, error)

type fileScanner struct {
	paths         []string
	excludedFiles []match.Matcher
	symlinks      bool

	log *logp.Logger
}

type fileWatcherConfig struct {
	Interval time.Duration
	Scanner  fileScannerConfig
}

type fileWatcher struct {
	interval time.Duration
	prev     map[string]os.FileInfo
	scanner  loginp.FSScanner
	log      *logp.Logger
	events   chan loginp.FSEvent
}

func newFileWatcher(paths []string, ns *common.ConfigNamespace) (loginp.FSWatcher, error) {
	if ns == nil {
		return newScannerWatcher(paths, nil)
	}

	watcherType := ns.Name()
	f, ok := watcherFactories[watcherType]
	if !ok {
		return nil, fmt.Errorf("no such file watcher: %s", watcherType)
	}

	return f(paths, ns.Config())
}

func newScannerWatcher(paths []string, c *common.Config) (loginp.FSWatcher, error) {
	config := defaultFileWatcherConfig()
	err := c.Unpack(&config)
	if err != nil {
		return nil, err
	}
	scanner, err := newFileScanner(paths, config.Scanner)
	if err != nil {
		return nil, err
	}
	return &fileWatcher{
		log:      logp.NewLogger(watcherDebugKey),
		interval: config.Interval,
		prev:     make(map[string]os.FileInfo, 0),
		scanner:  scanner,
		events:   make(chan loginp.FSEvent),
	}, nil
}

func defaultFileWatcherConfig() fileWatcherConfig {
	return fileWatcherConfig{
		Interval: 10 * time.Second,
		Scanner:  defaultFileScannerConfig(),
	}
}

func (w *fileWatcher) Run(ctx unison.Canceler) {
	defer close(w.events)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.watch(ctx)
		}
	}
}

// TODO add log messages
func (w *fileWatcher) watch(ctx unison.Canceler) {
	w.log.Info("Start next scan")

	paths := w.scanner.GetFiles()
	files := getKeys(paths)

	newFiles := make(map[string]os.FileInfo)

	for i := 0; i < len(paths); i++ {
		path := files[i]
		info := paths[path]

		prevInfo, ok := w.prev[path]
		if !ok {
			newFiles[path] = paths[path]
			continue
		}

		if prevInfo.ModTime() != info.ModTime() {
			w.log.Debug("updated file")
			select {
			case <-ctx.Done():
				return
			case w.events <- w.writeEvent(path, info):
			}
		}

		// delete from previous state, as we have more up to date info
		delete(w.prev, path)
	}

	// remaining files are in the prev map are the ones that are missing
	// either because they have been deleted or renamed
	for deletedPath, deletedInfo := range w.prev {
		for newPath, newInfo := range newFiles {
			if os.SameFile(deletedInfo, newInfo) {
				w.log.Debug("renamed file")
				select {
				case <-ctx.Done():
					return
				case w.events <- w.renamedEvent(deletedPath, newPath, newInfo):
					delete(newFiles, newPath)
					goto CHECK_NEXT_DELETED
				}
			}
		}

		w.log.Debug("deleted file")
		select {
		case <-ctx.Done():
			return
		case w.events <- w.deleteEvent(deletedPath, deletedInfo):
		}
	CHECK_NEXT_DELETED:
	}

	// remaining files in newFiles are new
	for path, info := range newFiles {
		w.log.Debug("new file")
		select {
		case <-ctx.Done():
			return
		case w.events <- w.createEvent(path, info):
		}

	}

	w.log.Debugf("Found %d paths", len(paths))
	w.prev = paths
}

func getKeys(paths map[string]os.FileInfo) []string {
	files := make([]string, 0)
	for file := range paths {
		files = append(files, file)
	}
	return files
}

func (w *fileWatcher) createEvent(path string, fi os.FileInfo) loginp.FSEvent {
	return loginp.FSEvent{Op: loginp.OpCreate, OldPath: "", NewPath: path, Info: fi}
}

func (w *fileWatcher) writeEvent(path string, fi os.FileInfo) loginp.FSEvent {
	return loginp.FSEvent{Op: loginp.OpWrite, OldPath: path, NewPath: path, Info: fi}
}

func (w *fileWatcher) renamedEvent(oldPath, path string, fi os.FileInfo) loginp.FSEvent {
	return loginp.FSEvent{Op: loginp.OpRename, OldPath: oldPath, NewPath: path, Info: fi}
}

func (w *fileWatcher) deleteEvent(path string, fi os.FileInfo) loginp.FSEvent {
	return loginp.FSEvent{Op: loginp.OpDelete, OldPath: path, NewPath: "", Info: fi}
}

func (w *fileWatcher) Event() loginp.FSEvent {
	return <-w.events
}

type fileScannerConfig struct {
	Paths         []string
	ExcludedFiles []match.Matcher
	Symlinks      bool
	RecursiveGlob bool
}

func defaultFileScannerConfig() fileScannerConfig {
	return fileScannerConfig{
		Symlinks:      false,
		RecursiveGlob: true,
	}
}

func newFileScanner(paths []string, cfg fileScannerConfig) (loginp.FSScanner, error) {
	fs := fileScanner{
		paths:         paths,
		excludedFiles: cfg.ExcludedFiles,
		symlinks:      cfg.Symlinks,
		log:           logp.NewLogger(scannerDebugKey),
	}
	err := fs.resolveRecursiveGlobs(cfg)
	if err != nil {
		return nil, err
	}
	err = fs.normalizeGlobPatterns()
	if err != nil {
		return nil, err
	}

	return &fs, nil
}

// resolveRecursiveGlobs expands `**` from the globs in multiple patterns
func (s *fileScanner) resolveRecursiveGlobs(c fileScannerConfig) error {
	if !c.RecursiveGlob {
		s.log.Debug("recursive glob disabled")
		return nil
	}

	s.log.Debug("recursive glob enabled")
	var paths []string
	for _, path := range s.paths {
		patterns, err := file.GlobPatterns(path, recursiveGlobDepth)
		if err != nil {
			return err
		}
		if len(patterns) > 1 {
			s.log.Debugf("%q expanded to %#v", path, patterns)
		}
		paths = append(paths, patterns...)
	}
	s.paths = paths
	return nil
}

// normalizeGlobPatterns calls `filepath.Abs` on all the globs from config
func (s *fileScanner) normalizeGlobPatterns() error {
	var paths []string
	for _, path := range s.paths {
		pathAbs, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get the absolute path for %s: %v", path, err)
		}
		paths = append(paths, pathAbs)
	}
	s.paths = paths
	return nil
}

func (s *fileScanner) GetFiles() map[string]os.FileInfo {
	paths := map[string]os.FileInfo{}

	for _, path := range s.paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			s.log.Errorf(msgGlobFailed, path, err)
			continue
		}

	OUTER:
		for _, file := range matches {
			if s.skippedFile(file) {
				continue
			}

			// If symlink is enabled, it is checked that original is not part of same input
			// If original is harvested by other input, states will potentially overwrite each other
			if s.isOriginalAndSymlinkConfigured(file, paths) {
				continue OUTER
			}

			fileInfo, err := os.Stat(file)
			if err != nil {
				s.log.Debug(msgStatFileFailed, file, err)
				continue
			}
			paths[file] = fileInfo
		}
	}

	return paths
}

func (s *fileScanner) skippedFile(file string) bool {
	if s.isFileExcluded(file) {
		s.log.Debugf(msgExcludedFile, file)
		return true
	}

	fileInfo, err := os.Lstat(file)
	if err != nil {
		s.log.Debugf(msgLstatFileFailed, file, err)
		return true
	}

	if fileInfo.IsDir() {
		s.log.Debugf(msgDirSkipped, file)
		return true
	}

	isSymlink := fileInfo.Mode()&os.ModeSymlink > 0
	if isSymlink && !s.symlinks {
		s.log.Debugf(msgSymlinkSkipped, file)
		return true
	}

	return false
}

func (s *fileScanner) isOriginalAndSymlinkConfigured(file string, paths map[string]os.FileInfo) bool {
	if s.symlinks {
		fileInfo, err := os.Stat(file)
		if err != nil {
			s.log.Debugf(msgStatFileFailed, file, err)
			return true
		}

		for _, finfo := range paths {
			if os.SameFile(finfo, fileInfo) {
				s.log.Info(msgSymlinkSameSkipped, file, finfo.Name())
				return true
			}
		}
	}
	return false
}

func (s *fileScanner) isFileExcluded(file string) bool {
	patterns := s.excludedFiles
	return len(patterns) > 0 && s.matchAny(patterns, file)
}

// MatchAny checks if the text matches any of the regular expressions
func (s *fileScanner) matchAny(matchers []match.Matcher, text string) bool {
	for _, m := range matchers {
		if m.MatchString(text) {
			return true
		}
	}
	return false
}
