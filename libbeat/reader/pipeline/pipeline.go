package pipeline

import (
	"io"

	"github.com/elastic/beats/libbeat/reader"
	"github.com/elastic/beats/libbeat/reader/multiline"
	"github.com/elastic/beats/libbeat/reader/readfile"
	"github.com/elastic/beats/libbeat/reader/readjson"
)

type Pipeline struct {
	config Config
	reader reader.Reader
}

type Config struct {
	File      readfile.Config
	JSON      *readjson.Config
	Multiline *multiline.Config
	// this is a hack to make docker JSON parsing part of the pipeline
	DockerJSON *struct {
		Stream   string `config:"stream"`
		Partial  bool   `config:"partial"`
		Format   string `config:"format"`
		CRIFlags bool   `config:"cri_flags"`
	}
	MaxBytes int
}

func New(r io.Reader, c Config) (Pipeline, error) {
	var reader reader.Reader
	var err error

	reader, err = readfile.NewEncodeReader(r, c.File)
	if err != nil {
		return Pipeline{}, err
	}

	return NewMessageReader(reader, c)
}

func NewMessageReader(reader reader.Reader, c Config) (Pipeline, error) {
	if c.DockerJSON != nil {
		// Docker json-file format, add custom parsing to the pipeline
		reader = readjson.New(reader, c.DockerJSON.Stream, c.DockerJSON.Partial, c.DockerJSON.Format, c.DockerJSON.CRIFlags)
	}

	if c.JSON != nil {
		reader = readjson.NewJSONReader(reader, c.JSON)
	}

	reader = readfile.NewStripNewline(reader, c.File.Terminator)

	if c.Multiline != nil {
		reader, err = multiline.New(reader, "\n", c.MaxBytes, c.Multiline)
		if err != nil {
			return Pipeline{}, err
		}
	}

	reader = readfile.NewLimitReader(reader, c.MaxBytes)

	return Pipeline{
		config: c,
		reader: reader,
	}, nil
}

func (p Pipeline) Next() (reader.Message, error) {
	return p.reader.Next()
}
