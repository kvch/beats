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

package filenotify

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/elastic/beats/filebeat/input/file"
	"github.com/elastic/beats/libbeat/common/match"
	"github.com/elastic/beats/libbeat/logp"
)

const (
	// Unchanged means the file is has not change since it was last read.
	Unchanged = iota
	// Append incites the file has new lines.
	Append
	// Truncated indicates that a file has been truncated.
	Truncated
	// Removed means that the file does not exist under its path.
	Removed
	// Created indicated that the file has been created.
	Created
	// Renamed is returned when the file being read is renamed.
	Renamed
	// Inactive means the file has not changed for the configured time in ignore_older.
	Inactive
	// ErrorReading means there was an error when checking its state.
	ErrorReading

	recursiveGlobDepth = 8
)

type Event struct {
	Change int
	Path   string
	Info   os.FileInfo
}

type FileSortInfo struct {
	info os.FileInfo
	path string
}

type Scanner struct {
	config Config

	globs         []string
	paths         []string
	previousState map[string]os.FileInfo
	timer         *time.Timer
	consumer      chan Event
	done          chan struct{}
}

func New(c Config, consumer chan Event, done chan struct{}) (*Scanner, error) {
	scanner := &Scanner{
		config:        c,
		previousState: make(map[string]os.FileInfo),
		timer:         time.NewTimer(c.ScanFrequency),
		consumer:      consumer,
		done:          done,
	}

	if err := scanner.resolveRecursiveGlobs(); err != nil {
		return nil, fmt.Errorf("Failed to resolve recursive globs in config: %v", err)
	}
	if err := scanner.normalizeGlobPatterns(); err != nil {
		return nil, fmt.Errorf("Failed to normalize globs patterns: %v", err)
	}

	return scanner, nil
}

// isFileExcluded checks if the given path should be excluded
func (s *Scanner) IsFileExcluded(file string) bool {
	patterns := s.config.ExcludeFiles
	return len(patterns) > 0 && matchAny(patterns, file)
}

// MatchAny checks if the text matches any of the regular expressions
func matchAny(matchers []match.Matcher, text string) bool {
	for _, m := range matchers {
		if m.MatchString(text) {
			return true
		}
	}
	return false
}

// Run scans files and notifies the input
func (s *Scanner) Run() {
	s.scan()
}

func (s *Scanner) scan() {
	logp.Debug("filenotify", "Start next scan")
	newFileInfos, err := s.getFileInfos()
	if err != nil {
		logp.Err("Error while scanning files: %v", err)
		return
	}

	logp.Info("PREV %v", s.previousState)
	logp.Info("NEXT %v", newFileInfos)
	s.notifyConsumer(newFileInfos)
	s.previousState = newFileInfos
}

func (s *Scanner) getFileInfos() (map[string]os.FileInfo, error) {
	var paths []string
	pathInfo := map[string]os.FileInfo{}

	go func() {
		select {
		case <-s.done:
			return
		}
	}()

	for _, g := range s.globs {
		matches, err := filepath.Glob(g)
		if err != nil {
			logp.Err("glob(%s) failed: %v", g, err)
			continue
		}

	OUTER:
		// Check any matched files to see if we need to start a harvester
		for _, file := range matches {

			// check if the file is in the exclude_files list
			if s.IsFileExcluded(file) {
				logp.Debug("filenotify", "Exclude file: %s", file)
				continue
			}

			// Fetch Lstat File info to detected also symlinks
			fileInfo, err := os.Lstat(file)
			if err != nil {
				logp.Debug("filenotify", "lstat(%s) failed: %s", file, err)
				continue
			}

			if fileInfo.IsDir() {
				logp.Debug("filenotify", "Skipping directory: %s", file)
				continue
			}

			isSymlink := fileInfo.Mode()&os.ModeSymlink > 0
			if isSymlink && !s.config.Symlinks {
				logp.Debug("filenotify", "File %s skipped as it is a symlink.", file)
				continue
			}

			// Fetch Stat file info which fetches the inode. In case of a symlink, the original inode is fetched
			fileInfo, err = os.Stat(file)
			if err != nil {
				logp.Debug("filenotify", "stat(%s) failed: %s", file, err)
				continue
			}

			// If symlink is enabled, it is checked that original is not part of same input
			// It original is harvested by other input, states will potentially overwrite each other
			if s.config.Symlinks {
				for _, finfo := range pathInfo {
					if os.SameFile(finfo, fileInfo) {
						logp.Info("Same file found as symlink and originap. Skipping file: %s", file)
						continue OUTER
					}
				}
			}
			logp.Info("EZT ADOM HOZZA %s", file)
			pathInfo[file] = fileInfo
			paths = append(paths, file)
		}
	}
	s.paths = paths
	return pathInfo, nil
}

func (s *Scanner) notifyConsumer(states map[string]os.FileInfo) {
	logp.Info("ATADTAM %d", len(states))
	logp.Info("ATADTAM %v", states)
	logp.Info("ATADTAM %v", s.paths)

	for _, p := range s.paths {
		f1, oneExists := s.previousState[p]
		f2, otherExists := states[p]
		logp.Info("ezek vannak %v %v", f1, f2)
		e := Event{
			Change: Unchanged,
			Path:   p,
			Info:   f2,
		}
		if oneExists && !otherExists {
			e.Change = Removed
		}
		if !oneExists && otherExists {
			e.Change = Created
		}

		if oneExists && otherExists {
			if f2.ModTime().Sub(f1.ModTime()) > 0 {
				if f1.Size() < f2.Size() {
					e.Change = Append
				} else {
					e.Change = Truncated
				}
			}
		}
		if time.Now().Sub(f2.ModTime()) > s.config.IgnoreOlder {
			e.Change = Inactive
		}
		s.consumer <- e
	}
}

// resolveRecursiveGlobs expands `**` from the globs in multiple patterns
func (s *Scanner) resolveRecursiveGlobs() error {
	if !s.config.RecursiveGlob {
		logp.Debug("filenotify", "recursive glob disabled")
		s.globs = s.config.Paths
		return nil
	}

	logp.Debug("filenotify", "recursive glob enabled")
	var paths []string
	for _, path := range s.config.Paths {
		patterns, err := file.GlobPatterns(path, recursiveGlobDepth)
		if err != nil {
			return err
		}
		if len(patterns) > 1 {
			logp.Debug("filenotify", "%q expanded to %#v", path, patterns)
		}
		paths = append(paths, patterns...)
	}
	s.globs = paths
	return nil
}

// normalizeGlobPatterns calls `filepath.Abs` on all the globs from config
func (s *Scanner) normalizeGlobPatterns() error {
	var paths []string
	for _, path := range s.globs {
		pathAbs, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("Failed to get the absolute path for %s: %v", path, err)
		}
		paths = append(paths, pathAbs)
	}
	s.globs = paths
	return nil
}

func getSortInfos(paths map[string]os.FileInfo) []FileSortInfo {
	sortInfos := make([]FileSortInfo, 0, len(paths))
	for path, info := range paths {
		sortInfo := FileSortInfo{info: info, path: path}
		sortInfos = append(sortInfos, sortInfo)
	}

	return sortInfos
}

func getSortedFiles(scanOrder string, scanSort string, sortInfos []FileSortInfo) ([]FileSortInfo, error) {
	var sortFunc func(i, j int) bool
	switch scanSort {
	case "modtime":
		switch scanOrder {
		case "asc":
			sortFunc = func(i, j int) bool {
				return sortInfos[i].info.ModTime().Before(sortInfos[j].info.ModTime())
			}
		case "desc":
			sortFunc = func(i, j int) bool {
				return sortInfos[i].info.ModTime().After(sortInfos[j].info.ModTime())
			}
		default:
			return nil, fmt.Errorf("Unexpected value for scan.order: %v", scanOrder)
		}
	case "filename":
		switch scanOrder {
		case "asc":
			sortFunc = func(i, j int) bool {
				return strings.Compare(sortInfos[i].info.Name(), sortInfos[j].info.Name()) < 0
			}
		case "desc":
			sortFunc = func(i, j int) bool {
				return strings.Compare(sortInfos[i].info.Name(), sortInfos[j].info.Name()) > 0
			}
		default:
			return nil, fmt.Errorf("Unexpected value for scan.order: %v", scanOrder)
		}
	default:
		return nil, fmt.Errorf("Unexpected value for scan.sort: %v", scanSort)
	}

	if sortFunc != nil {
		sort.Slice(sortInfos, sortFunc)
	}

	return sortInfos, nil
}
