package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/elastic/beats/filebeat/fileset"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/cmd"
)

func buildModulesManager(beat *beat.Beat) (cmd.ModulesManager, error) {
	config := beat.BeatConfig

	glob, err := config.String("config.modules.path", -1)
	if err != nil {
		return nil, errors.Errorf("modules management requires 'filebeat.config.modules.path' setting")
	}

	if !strings.HasSuffix(glob, "*.yml") {
		return nil, errors.Errorf("wrong settings for config.modules.path, it is expected to end with *.yml. Got: %s", glob)
	}

	modulesManager, err := cfgfile.NewGlobManager(glob, ".yml", ".disabled")
	if err != nil {
		return nil, errors.Wrap(err, "initialization error")
	}

	return modulesManager, nil
}

func genGenerateCmd() *cobra.Command {
	genCmd := cobra.Command{
		Use:   "generate ACTION MODULE FILESET",
		Short: "Generate new Filebeat module: fileset, fields",
		Long: `Generate new fileset for Filebeat. Possible actions: fileset, fields.
 * "fileset" creates directories and config files for you new fileset.
 * "field" generates fields.yml based on pipeline.json. Must be called after pipeline.json is finished.`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				fmt.Fprintf(os.Stderr, "Three arguments are required: action module fileset, got: %d\n", len(args))
				os.Exit(1)

			}

			if args[0] == "fields" {
				err := fileset.GenerateFieldsYml(args[1], args[2])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to generate fields.yml for %s/%s: %v\n", args[1], args[2], err)
					os.Exit(1)
				}
			} else if args[0] == "fileset" {
				err := fileset.Generate(args[1], args[2])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to generate fileset %s/%s: %v\n", args[1], args[2], err)
					os.Exit(1)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Unknown action: %s\n", args[0])
				os.Exit(1)
			}
		},
	}
	genCmd.Flags().BoolVarP(&fileset.Nodoc, "nodoc", "n", false, "Add documentation fields for generated fields.yml")
	return &genCmd
}
