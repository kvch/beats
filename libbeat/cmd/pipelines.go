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
	pipelinesURL = "/_ingest/pipeline/%s"
)

// GenModulesCmd initializes a command to manage a modules.d folder, it offers
// list, enable and siable actions
func genPipelinesCmd(name, version string) *cobra.Command {
	pipelinesCmd := cobra.Command{
		Use:   "pipelines",
		Short: "Manage Ingest pipelines",
	}

	pipelinesCmd.AddCommand(genListPipelinesCmd(name, version))
	pipelinesCmd.AddCommand(genDeletePipelinesCmd(name, version))

	return &pipelinesCmd
}

func genListPipelinesCmd(name, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "list [id]",
		Short: "List pipelines",
		Run: func(cmd *cobra.Command, args []string) {
			esClient, err := setupElasticsearchClient(name, version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				os.Exit(1)
			}

			reqURL := ""
			if len(args) == 0 {
				reqURL = fmt.Sprintf(pipelinesURL, "")
			} else if len(args) == 1 {
				reqURL = fmt.Sprintf(pipelinesURL, args[0])
			} else {
				fmt.Fprintf(os.Stderr, "Too many parameters\n")
				os.Exit(2)
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

			err = printResponse(r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				os.Exit(3)
			}
		},
	}
}

func genDeletePipelinesCmd(name, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete pipelines",
		Run: func(cmd *cobra.Command, args []string) {
			esClient, err := setupElasticsearchClient(name, version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				os.Exit(1)
			}

			reqURL := ""
			if len(args) == 1 {
				reqURL = fmt.Sprintf(pipelinesURL, args[0])
			} else {
				fmt.Fprintf(os.Stderr, "One parameter is required\n")
				os.Exit(2)
			}

			code, r, err := esClient.Request("DELETE", reqURL, "", nil, nil)
			if code == 404 {
				fmt.Println("No such pipeline")
				os.Exit(0)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error during request to Elasticsearch: %v\n", err)
				os.Exit(1)
			}

			err = printResponse(r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				os.Exit(3)
			}
		},
	}
}

func setupElasticsearchClient(name, version string) (*elasticsearch.Client, error) {
	b, err := instance.NewBeat(name, "", version)
	if err != nil {
		return nil, fmt.Errorf("Error initializing beat: %s\n", err)
	}

	if err = b.Init(); err != nil {
		return nil, fmt.Errorf("Error initializing beat: %s\n", err)
	}

	if b.Config.Output.Name() != "elasticsearch" {
		return nil, fmt.Errorf("Error conencting to Elasticsearch: missing configuration from file\n")
	}

	esConfig := b.Config.Output.Config()
	esClient, err := elasticsearch.NewConnectedClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("Error conencting to Elasticsearch: %v\n", err)
	}
	return esClient, nil
}

func printResponse(response []byte) error {
	var resp common.MapStr
	buf := bytes.NewBuffer(response)
	err := json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		return fmt.Errorf("Error during parsing response from Elasticsearch: %v\n", err)
	}

	fmt.Println(resp.StringToPrint())

	return nil
}
