package beater

import (
	"fmt"
	"sync"

	"github.com/elastic/beats/journalbeat/input"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/elastic/beats/journalbeat/config"
)

type Journalbeat struct {
	inputs []*input.Input
	done   chan struct{}
	config config.Config

	client   beat.Client
	pipeline beat.Pipeline
}

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	done := make(chan struct{})

	var inputs []*input.Input
	for _, c := range config.Inputs {
		client, err := b.Publisher.Connect()
		if err != nil {
			return nil, err
		}

		i := input.New(c, client, done)
		// TODO
		if i == nil {
			continue
		}
		inputs = append(inputs, i)
	}

	bt := &Journalbeat{
		inputs: inputs,
		done:   done,
		config: config,
	}
	return bt, nil
}

func (bt *Journalbeat) Run(b *beat.Beat) error {
	logp.Info("journalbeat is running! Hit CTRL-C to stop it.")
	defer logp.Info("journalbeat is stopping")

	var wg sync.WaitGroup
	for _, i := range bt.inputs {
		wg.Add(1)
		go bt.runInput(i, &wg)
	}
	wg.Wait()

	return nil
}

func (bt *Journalbeat) runInput(i *input.Input, wg *sync.WaitGroup) {
	defer wg.Done()
	i.Run()
}

func (bt *Journalbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
