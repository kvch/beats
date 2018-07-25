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

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/elastic/beats/libbeat/generator/fields"
)

func main() {
	esBeatsPath := flag.String("es_beats_path", "..", "Path to elastic/beats")
	beatPath := flag.String("beat_path", ".", "Path to your Beat")
	flag.Parse()

	beatFieldsPaths := flag.Args()
	name := filepath.Base(*beatPath)

	if *beatPath == "" {
		fmt.Fprintf(os.Stderr, "beat_path cannot be empty")
		os.Exit(1)
	}

	esBeats, err := os.Open(*esBeatsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening elastic/beats: %+v\n", err)
		os.Exit(1)
	}
	beat, err := os.Open(*beatPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening target Beat: %+v\n", err)
		os.Exit(1)
	}
	esBeatsInfo, err := esBeats.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting file info of elastic/beat: %+v\n", err)
		os.Exit(1)
	}
	beatInfo, err := beat.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting file info of target Beat: %+v\n", err)
		os.Exit(1)
	}

	if len(beatFieldsPaths) == 0 && os.SameFile(esBeatsInfo, beatInfo) {
		fmt.Println("No field files to collect")
		return
	}

	var fieldsFiles []*fields.YmlFile
	for _, fieldsFilePath := range beatFieldsPaths {
		pathToModules := filepath.Join(*beatPath, fieldsFilePath)

		fieldsFile, err := fields.CollectModuleFiles(pathToModules)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot collect fields.yml files: %+v\n", err)
			os.Exit(2)
		}

		fieldsFiles = append(fieldsFiles, fieldsFile...)
	}

	err = fields.Generate(*esBeatsPath, *beatPath, fieldsFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot generate global fields.yml file for %s: %+v\n", name, err)
		os.Exit(3)
	}

	fmt.Printf("Generated fields.yml for %s\n", name)
}
