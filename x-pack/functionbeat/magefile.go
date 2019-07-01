// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/magefile/mage/mg"

	devtools "github.com/elastic/beats/dev-tools/mage"
	functionbeat "github.com/elastic/beats/x-pack/functionbeat/scripts/mage"
)

var (
	availableProviders = []string{
		"aws",
	}
	selectedProviders []string
)

func init() {
	devtools.BeatDescription = "Functionbeat is a beat implementation for a serverless architecture."
	devtools.BeatLicense = "Elastic License"
	selectedProviders = getConfiguredProviders()
}

// Build builds the Beat binary.
func Build() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	params := devtools.DefaultBuildArgs()
	for _, provider := range selectedProviders {
		params.Name = devtools.BeatName + "-" + provider
		err = os.Chdir(workingDir + "/" + provider)
		if err != nil {
			return err
		}

		err = devtools.Build(params)
		if err != nil {
			return err
		}
	}
	return nil
}

// GolangCrossBuild build the Beat binary inside of the golang-builder.
// Do not use directly, use crossBuild instead.
func GolangCrossBuild() error {
	params := devtools.DefaultGolangCrossBuildArgs()
	params.Name = "functionbeat-" + params.Name
	return devtools.GolangCrossBuild(params)
}

// BuildGoDaemon builds the go-daemon binary (use crossBuildGoDaemon).
func BuildGoDaemon() error {
	return devtools.BuildGoDaemon()
}

// CrossBuild cross-builds the beat for all target platforms.
func CrossBuild() error {
	for _, provider := range selectedProviders {
		err := devtools.CrossBuild(devtools.AddPlatforms("linux/amd64"), devtools.InDir("x-pack", "functionbeat", provider))
		if err != nil {
			return err
		}
	}
	return nil
}

// CrossBuildGoDaemon cross-builds the go-daemon binary using Docker.
func CrossBuildGoDaemon() error {
	return devtools.CrossBuildGoDaemon()
}

// Clean cleans all generated files and build artifacts.
func Clean() error {
	return devtools.Clean()
}

// Package packages the Beat for distribution.
// Use SNAPSHOT=true to build snapshots.
// Use PLATFORMS to control the target platforms.
func Package() {
	start := time.Now()
	defer func() { fmt.Println("package ran for", time.Since(start)) }()

	mg.Deps(Update)
	mg.Deps(CrossBuild, CrossBuildGoDaemon)
	for _, provider := range selectedProviders {
		devtools.MustUsePackaging("functionbeat", "x-pack/functionbeat/dev-tools/packaging/packages.yml")
		for _, args := range devtools.Packages {
			args.Spec.ExtraVar("Provider", provider)
		}

		mg.SerialDeps(devtools.Package, TestPackages)
	}
}

// TestPackages tests the generated packages (i.e. file modes, owners, groups).
func TestPackages() error {
	return devtools.TestPackages()
}

// Update updates the generated files (aka make update).
func Update() {
	mg.SerialDeps(Fields, Config, includeFields, docs)
}

// GoTestUnit executes the Go unit tests.
// Use TEST_COVERAGE=true to enable code coverage profiling.
// Use RACE_DETECTOR=true to enable the race detector.
func GoTestUnit(ctx context.Context) error {
	return devtools.GoTest(ctx, devtools.DefaultGoTestUnitArgs())
}

// GoTestIntegration executes the Go integration tests.
// Use TEST_COVERAGE=true to enable code coverage profiling.
// Use RACE_DETECTOR=true to enable the race detector.
func GoTestIntegration(ctx context.Context) error {
	return devtools.GoTest(ctx, devtools.DefaultGoTestIntegrationArgs())
}

// Config generates both the short and reference configs.
func Config() error {
	for _, provider := range selectedProviders {
		devtools.BeatIndexPrefix += "-" + provider
		err := devtools.Config(devtools.ShortConfigType|devtools.ReferenceConfigType, functionbeat.XPackConfigFileParams(provider), provider)
		if err != nil {
			return err
		}
	}
	return nil
}

// Fields generates a fields.yml for the Beat.
func Fields() error {
	for _, provider := range selectedProviders {
		output := filepath.Join(devtools.CWD(), provider, "fields.yml")
		err := devtools.GenerateFieldsYAMLTo(output)
		if err != nil {
			return err
		}
	}
	return nil
}

func includeFields() error {
	fnBeatDir := devtools.CWD()
	for _, provider := range selectedProviders {
		err := os.Chdir(filepath.Join(fnBeatDir, provider))
		if err != nil {
			return err
		}
		output := filepath.Join(fnBeatDir, provider, "include", "fields.go")
		err = devtools.GenerateFieldsGoWithName(devtools.BeatName+"-"+provider, "fields.yml", output)
		if err != nil {
			return err
		}
	}
	os.Chdir(fnBeatDir)
	return nil
}

func docs() error {
	for _, provider := range selectedProviders {
		fieldsYml := filepath.Join(devtools.CWD(), provider, "fields.yml")
		err := devtools.Docs.FieldDocs(fieldsYml)
		if err != nil {
			return err
		}
	}
	return nil
}

// IntegTest executes integration tests (it uses Docker to run the tests).
func IntegTest() {
	devtools.AddIntegTestUsage()
	defer devtools.StopIntegTestEnv()
	mg.SerialDeps(GoIntegTest, PythonIntegTest)
}

// GoIntegTest executes the Go integration tests.
// Use TEST_COVERAGE=true to enable code coverage profiling.
// Use RACE_DETECTOR=true to enable the race detector.
func GoIntegTest(ctx context.Context) error {
	return devtools.RunIntegTest("goIntegTest", func() error {
		return devtools.GoTest(ctx, devtools.DefaultGoTestIntegrationArgs())
	})
}

// PythonUnitTest executes the python system tests.
func PythonUnitTest() error {
	mg.Deps(devtools.BuildSystemTestBinary)
	return devtools.PythonNoseTest(devtools.DefaultPythonTestUnitArgs())
}

// PythonIntegTest executes the python system tests in the integration environment (Docker).
func PythonIntegTest(ctx context.Context) error {
	if !devtools.IsInIntegTestEnv() {
		mg.Deps(Fields)
	}
	return devtools.RunIntegTest("pythonIntegTest", func() error {
		mg.Deps(devtools.BuildSystemTestBinary)
		args := devtools.DefaultPythonTestIntegrationArgs()

		workingDir := devtools.CWD()
		for _, provider := range selectedProviders {
			args.Env = map[string]string{
				"CURRENT_PROVIDER": "aws",
			}
			err := os.Chdir(workingDir + "/" + provider)
			if err != nil {
				return err
			}
			err = devtools.PythonNoseTest(args)
			if err != nil {
				return err
			}
		}
		return os.Chdir(workingDir)
	})
}

// BuildSystemTestBinary build a binary for testing that is instrumented for
// testing and measuring code coverage. The binary is only instrumented for
// coverage when TEST_COVERAGE=true (default is false).
func BuildSystemTestBinary() error {
	workingDir := devtools.CWD()
	for _, provider := range selectedProviders {
		err := os.Chdir(workingDir + "/" + provider)
		if err != nil {
			return err
		}
		err = devtools.BuildSystemTestBinary(devtools.TestBinaryArgs{Name: devtools.BeatName + "-" + provider})
		if err != nil {
			return err
		}
	}
	return os.Chdir(workingDir)
}

func getConfiguredProviders() []string {
	providers := os.Getenv("PROVIDERS")
	if len(providers) == 0 {
		return availableProviders
	}

	return strings.Split(providers, ",")
}
