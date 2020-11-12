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
	"github.com/elastic/go-concert/unison"
)

const (
	prospectorDebugKey = "file_prospector"
)

// fileProspector implements the Prospector interface.
// It contains a file scanner which returns file system events.
// The FS events then trigger either new Harvester runs or updates
// the statestore.
type fileProspector struct {
	filewatcher  loginp.FSWatcher
	identifier   fileIdentifier
	ignoreOlder  time.Duration
	cleanRemoved bool
}

func newFileProspector(
	paths []string,
	ignoreOlder time.Duration,
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
	}, nil
}

func (p *fileProspector) Init(cleaner loginp.ProspectorCleaner) error {
	if p.cleanRemoved {
		cleaner.CleanIf(func(key string, v loginp.Value) bool {
			var fm fileMeta
			err := v.UnpackCursorMeta(&fm)
			if err != nil {
				// remove faulty entries
				return true
			}

			_, err = os.Stat(fm.Source)
			if err != nil {
				return true
			}
			return false
		})
	}

	identifierName := p.identifier.Name()
	cleaner.UpdateIdentifiers(func(key string, v loginp.Value) (bool, string, interface{}) {
		var fm fileMeta
		err := v.UnpackCursorMeta(&fm)
		if err != nil {
			return false, "", nil
		}
		if fm.IdentifierName != identifierName {
			fi, err := os.Stat(fm.Source)
			if err != nil {
				return false, "", fm
			}
			newKey := p.identifier.GetSource(loginp.FSEvent{NewPath: fm.Source, Info: fi}).Name()
			fm.IdentifierName = identifierName
			return true, newKey, fm
		}
		return false, "", nil
	})

	return nil
}

// Run starts the fileProspector which accepts FS events from a file watcher.
func (p *fileProspector) Run(ctx input.Context, s loginp.StateMetadataUpdater, hg loginp.HarvesterGroup) {
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
			case loginp.OpCreate, loginp.OpWrite:
				if fe.Op == loginp.OpCreate {
					log.Debugf("A new file %s has been found", fe.NewPath)
				} else if fe.Op == loginp.OpWrite {
					log.Debugf("File %s has been updated", fe.NewPath)
				}

				if p.ignoreOlder > 0 {
					now := time.Now()
					if now.Sub(fe.Info.ModTime()) > p.ignoreOlder {
						log.Debugf("Ignore file because ignore_older reached. File %s", fe.NewPath)
						break
					}
				}

				hg.Start(ctx, src)

			case loginp.OpDelete:
				log.Debugf("File %s has been removed", fe.OldPath)

				if p.cleanRemoved {
					log.Debugf("Remove state for file as file removed: %s", fe.OldPath)

					err := s.Remove(src.Name())
					if err != nil {
						log.Errorf("Error while removing state from statestore: %v", err)
					}
				}

				// TODO close_removed

			case loginp.OpRename:
				log.Debugf("File %s has been renamed to %s", fe.OldPath, fe.NewPath)

				// if file_identity is based on path, the current reader has to be cancelled
				// and a new one has to start.
				if p.identifier.Name() == "path" {
					prevSrc := p.identifier.GetSource(loginp.FSEvent{NewPath: fe.OldPath})
					hg.Stop(prevSrc)

					log.Debugf("Remove state for file as file renamed and path file_identity is configured: %s", fe.OldPath)
					err := s.Remove(prevSrc.Name())
					if err != nil {
						log.Errorf("Error while removing old state of renamed file (%s): %v", fe.OldPath, err)
					}

					hg.Start(ctx, src)
				} else {
					// update file metadata as the path has changed
					id := src.Name()
					var meta fileMeta
					err := s.FindCursorMeta(id, meta)
					if err != nil {
						log.Errorf("Error while getting cursor meta data of entry %s: %v", id, err)
					}
					s.UpdateMetadata(id, fileMeta{Source: src.newPath, IdentifierName: meta.IdentifierName})

					// TODO close_renamed
				}

			default:
				log.Error("Unkown return value %v", fe.Op)
			}
		}
		return nil
	})

	errs := tg.Wait()
	if len(errs) > 0 {
		log.Error("%s", sderr.WrapAll(errs, "running prospector failed"))
	}
}

func (p *fileProspector) Test() error {
	panic("TODO: implement me")
}
