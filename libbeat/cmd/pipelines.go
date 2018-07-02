package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/elastic/beats/libbeat/cmd/instance"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/outputs/elasticsearch"
)

const (
	listPipelines = "/_ingest/pipeline/%s"
)

// GenModulesCmd initializes a command to manage a modules.d folder, it offers
// list, enable and siable actions
func genPipelinesCmd(name, version string) *cobra.Command {
	pipelinesCmd := cobra.Command{
		Use:   "pipelines",
		Short: "Manage Ingest pipelines",
	}

	pipelinesCmd.AddCommand(genListPipelinesCmd(name, version))

	return &pipelinesCmd
}

func genListPipelinesCmd(name, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "list [id]",
		Short: "List pipelines",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := instance.NewBeat(name, "", version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing beat: %s\n", err)
				os.Exit(1)
			}

			if err = b.Init(); err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing beat: %s\n", err)
				os.Exit(1)
			}

			if b.Config.Output.Name() != "elasticsearch" {
				fmt.Fprintf(os.Stderr, "Error conencting to Elasticsearch: missing configuration from file\n")
				os.Exit(1)
			}

			reqURL := ""
			if len(args) == 0 {
				reqURL = fmt.Sprintf(listPipelines, "")
			} else if len(args) == 1 {
				reqURL = fmt.Sprintf(listPipelines, args[0])
			} else {
				fmt.Fprintf(os.Stderr, "Too many parameters\n")
				os.Exit(2)
			}

			esConfig := b.Config.Output.Config()
			esClient, err := elasticsearch.NewConnectedClient(esConfig)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error conencting to Elasticsearch: %v\n", err)
				os.Exit(1)
			}
			code, r, err := esClient.Request("GET", reqURL, "", nil, nil)
			if code == 404 {
				fmt.Println("No such pipeline")
				os.Exit(0)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error during request to Elasticsearch: %v\n", err)
				os.Exit(1)
			}

			var resp common.MapStr
			buf := bytes.NewBuffer(r)
			err = json.Unmarshal(buf.Bytes(), &resp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error during parsing response from Elasticsearch: %v\n", err)
				os.Exit(3)
			}
			fmt.Println(resp.StringToPrint())
		},
	}
}
