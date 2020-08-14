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
	"os"
	"time"

	"github.com/urso/sderr"

	loginp "github.com/elastic/beats/v7/filebeat/input/filestream/internal/input-logfile"
	input "github.com/elastic/beats/v7/filebeat/input/v2"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/statestore"
	"github.com/elastic/go-concert/unison"
)

const (
	prospectorDebugKey = "file_prospector"
)

type fileProspector struct {
	filewatcher  loginp.FSWatcher
	identifier   fileIdentifier
	ignoreOlder  time.Duration
	cleanRemoved bool // TODO
	monitor      *activeFileMonitor
}

func newFileProspector(
	paths []string,
	ignoreOlder time.Duration,
	closerConfig stateChangeCloserConfig,
	fileWatcherNs, identifierNs *common.ConfigNamespace,
) (loginp.Prospector, error) {

	filewatcher, err := newFileWatcher(paths, fileWatcherNs)
	if err != nil {
		return nil, err
	}

	identifier, err := newFileIdentifier(identifierNs)
	if err != nil {
		return nil, err
	}

	return &fileProspector{
		filewatcher:  filewatcher,
		identifier:   identifier,
		ignoreOlder:  ignoreOlder,
		cleanRemoved: true,
		monitor:      newActiveFileMonitor(closerConfig),
	}, nil
}

func (p *fileProspector) Run(ctx input.Context, s *statestore.Store, hg *loginp.HarvesterGroup) {
	log := ctx.Logger.With("prospector", prospectorDebugKey)
	log.Debug("Starting prospector")
	defer log.Debug("Prospector has stopped")

	var tg unison.MultiErrGroup

	tg.Go(func() error {
		p.filewatcher.Run(ctx.Cancelation)
		return nil
	})

	tg.Go(func() error {
		for ctx.Cancelation.Err() == nil {
			fe := p.filewatcher.Event()

			if fe.Op == loginp.OpDone {
				return nil
			}

			src := p.identifier.GetSource(fe)
			switch fe.Op {
			case loginp.OpCreate:
				log.Debugf("A new file %s has been found", fe.NewPath)

				if p.ignoreOlder > 0 {
					now := time.Now()
					if now.Sub(fe.Info.ModTime()) > p.ignoreOlder {
						log.Debugf("Ignore file because ignore_older reached. File %s", fe.NewPath)
						break
					}
				}

				p.startReading(ctx, hg, src, fe.NewPath)

			case loginp.OpWrite:
				log.Debugf("File %s has been updated", fe.NewPath)

				p.startReading(ctx, hg, src, fe.NewPath)

			case loginp.OpDelete:
				log.Debugf("File %s has been removed", fe.OldPath)

				// TODO when to clean up deleted files between Filebeat runs?
				if p.cleanRemoved {
					log.Debugf("Remove state for file as file removed: %s", fe.OldPath)

					err := s.Remove(src.Name())
					if err != nil {
						log.Errorf("Error while removing state from statestore: %v", err)
					}
				}

			case loginp.OpRename:
				log.Debugf("File %s has been renamed to %s", fe.OldPath, fe.NewPath)
				// TODO update state information in the store

			default:
				log.Error("Unkown return value %v", fe.Op)
			}
		}
		return nil
	})

	tg.Go(func() error {
		p.monitor.run(ctx.Cancelation, hg)
		return nil
	})

	errs := tg.Wait()
	if len(errs) > 0 {
		log.Error("%s", sderr.WrapAll(errs, "running prospector failed"))
	}
}

func (p *fileProspector) startReading(ctx input.Context, hg *loginp.HarvesterGroup, s fileSource, path string) {
	added := p.monitor.addFile(path, s)
	if added {
		hg.Run(ctx, s)
	}
}

// isSameFile checks if the given File path corresponds with the FileInfo given
// It is used to check if the file has been renamed.
func isSameFile(path string, info os.FileInfo) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return os.SameFile(fileInfo, info), nil
}

func (p *fileProspector) Test() error {
	return nil
}
