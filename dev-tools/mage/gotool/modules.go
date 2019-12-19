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

package gotool

type goMod func(opts ...ArgOpt) error

// Mod is the command go mod.
var Mod goMod = runGoMod

// Tidy cleans the go.mod file
func (goMod) Tidy() error {
	return runGoMod(posArg("tidy"))
}

// Vendor downloads and copies dependencies under the folder vendor.
func (goMod) Vendor() error {
	return runGoMod(posArg("vendor"))
}

// Verify check that deps have the expected content.
func (goMod) Verify() error {
	return runGoMod(posArg("verify"))
}

func runGoMod(opts ...ArgOpt) error {
	args := buildArgs(opts)
	return runVGo("mod", args)
}
