package fileset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	pipelinePath  = "../module/%s/%s/ingest/pipeline.json"
	fieldsYmlPath = "../module/%s/%s/_meta/fields.yml"
)

var (
	types = map[string]string{
		"group":           "group",
		"DATA":            "text",
		"GREEDYDATA":      "text",
		"GREEDYMULTILINE": "text",
		"HOSTNAME":        "keyword",
		"IPHOST":          "keyword",
		"IPORHOST":        "keyword",
		"LOGLEVEL":        "keyword",
		"MULTILINEQUERY":  "text",
		"NUMBER":          "long",
		"POSINT":          "long",
		"SYSLOGHOST":      "keyword",
		"SYSLOGTIMESTAMP": "text",
		"TIMESTAMP":       "text",
		"USERNAME":        "keyword",
		"WORD":            "keyword",
	}
	// Nodoc specifies if the generated field.yml is includes documentation fields e.g. description, example
	Nodoc bool
)

type pipeline struct {
	Description string                   `json:"description"`
	Processors  []map[string]interface{} `json:"processors"`
	OnFailure   interface{}              `json:"on_failure"`
}

type field struct {
	Type     string
	Elements []string
}

type fieldYml struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description,omitempty"`
	Example     string      `yaml:"example,omitempty"`
	Type        string      `yaml:"type,omitempty"`
	Fields      []*fieldYml `yaml:"fields,omitempty"`
}

func newFieldYml(name, typeName string, noDoc bool) *fieldYml {
	if noDoc {
		return &fieldYml{
			Name: name,
			Type: typeName,
		}
	}

	return &fieldYml{
		Name:        name,
		Type:        typeName,
		Description: "Please add description",
		Example:     "Please add example",
	}
}

func newField(lp string) field {
	lp = lp[1 : len(lp)-1]
	ee := strings.Split(lp, ":")
	e := strings.Split(ee[1], ".")
	return field{
		Type:     ee[0],
		Elements: e,
	}
}

func readPipeline(module, fileset string) (*pipeline, error) {
	pp := fmt.Sprintf(pipelinePath, module, fileset)
	r, err := ioutil.ReadFile(pp)
	if err != nil {
		return nil, err
	}

	var p pipeline
	err = json.Unmarshal(r, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func addNewField(fs []field, f field) []field {
	for _, ff := range fs {
		if reflect.DeepEqual(ff, f) {
			return fs
		}
	}
	return append(fs, f)
}

func getElementsFromPatterns(patterns []string) ([]field, error) {
	r, err := regexp.Compile("{[\\.\\w\\:]*}")
	if err != nil {
		return nil, err
	}

	fs := make([]field, 0)
	for _, lp := range patterns {
		pp := r.FindAllString(lp, -1)
		for _, p := range pp {
			f := newField(p)
			fs = addNewField(fs, f)
		}

	}
	return fs, nil
}

func accumulatePatterns(grok interface{}) ([]string, error) {
	for k, v := range grok.(map[string]interface{}) {
		if k == "patterns" {
			vs := v.([]interface{})
			p := make([]string, 0)
			for _, s := range vs {
				p = append(p, s.(string))
			}
			return p, nil
		}
	}
	return nil, fmt.Errorf("No patterns in pipeline")
}

func accumulateRemoveFields(remove interface{}, out []string) []string {
	for k, v := range remove.(map[string]interface{}) {
		if k == "field" {
			vs := v.(string)
			return append(out, vs)
		}
	}
	return out
}

func accumulateRenameFields(rename interface{}, out map[string]string) map[string]string {
	var from, to string
	for k, v := range rename.(map[string]interface{}) {
		if k == "field" {
			from = v.(string)
		}
		if k == "target_field" {
			to = v.(string)
		}
	}
	out[from] = to
	return out
}

type processors struct {
	patterns []string
	remove   []string
	rename   map[string]string
}

func (p *processors) processFields() ([]field, error) {
	f, err := getElementsFromPatterns(p.patterns)
	if err != nil {
		return nil, err
	}

	for i, ff := range f {
		fs := strings.Join(ff.Elements, ".")
		for _, rm := range p.remove {
			if fs == rm {
				f = append(f[:i], f[i+1:]...)
			}
		}
		for k, mv := range p.rename {
			if k == fs {
				ff.Elements = strings.Split(mv, ".")
			}
		}
		f[i] = ff
	}
	return f, nil
}

func getProcessors(p []map[string]interface{}) (*processors, error) {
	patterns := make([]string, 0)
	rmFields := make([]string, 0)
	mvFields := make(map[string]string)
	var err error

	for _, e := range p {
		if ee, ok := e["grok"]; ok {
			patterns, err = accumulatePatterns(ee)
			if err != nil {
				return nil, err
			}
		}
		if rm, ok := e["remove"]; ok {
			rmFields = accumulateRemoveFields(rm, rmFields)
		}
		if mv, ok := e["rename"]; ok {
			mvFields = accumulateRenameFields(mv, mvFields)
		}
	}

	if patterns == nil {
		return nil, fmt.Errorf("No patterns in pipeline")
	}

	return &processors{
		patterns: patterns,
		remove:   rmFields,
		rename:   mvFields,
	}, nil
}

func getFieldByName(f []*fieldYml, name string) *fieldYml {
	for _, ff := range f {
		if ff.Name == name {
			return ff
		}
	}
	return nil
}

func insertLastField(f []*fieldYml, name, typeName string, noDoc bool) []*fieldYml {
	ff := getFieldByName(f, name)
	if ff != nil {
		return f
	}

	nf := newFieldYml(name, types[typeName], noDoc)
	return append(f, nf)
}

func insertGroup(out []*fieldYml, field field, index, count int, noDoc bool) []*fieldYml {
	g := getFieldByName(out, field.Elements[index])
	if g != nil {
		g.Fields = generateField(g.Fields, field, index+1, count, noDoc)
		return out
	}

	groupFields := make([]*fieldYml, 0)
	groupFields = generateField(groupFields, field, index+1, count, noDoc)
	group := newFieldYml(field.Elements[index], "group", noDoc)
	group.Fields = groupFields
	return append(out, group)
}

func generateField(out []*fieldYml, field field, index, count int, noDoc bool) []*fieldYml {
	if index+1 == count {
		return insertLastField(out, field.Elements[index], field.Type, noDoc)
	}
	return insertGroup(out, field, index, count, noDoc)
}

func generateFields(f []field, noDoc bool) []*fieldYml {
	out := make([]*fieldYml, 0)
	for _, ff := range f {
		out = generateField(out, ff, 1, len(ff.Elements), noDoc)
	}
	return out
}

func (p *pipeline) toFieldsYml(noDoc bool) ([]byte, error) {
	pr, err := getProcessors(p.Processors)
	if err != nil {
		return nil, err
	}

	var fs []field
	fs, err = pr.processFields()
	if err != nil {
		return nil, err
	}

	f := generateFields(fs, noDoc)
	var d []byte
	d, err = yaml.Marshal(&f)

	return d, nil
}

func writeFieldsYml(module, fileset string, f []byte) error {
	p := fmt.Sprintf(fieldsYmlPath, module, fileset)
	err := ioutil.WriteFile(p, f, 0664)
	if err != nil {
		return err
	}
	return nil
}

// GenerateFieldsYml creates a fields.yml file for fileset based on the existing pipeline.json
func GenerateFieldsYml(module, fileset string) error {
	p, err := readPipeline(module, fileset)
	if err != nil {
		return err
	}

	var d []byte
	d, err = p.toFieldsYml(Nodoc)
	if err != nil {
		return err
	}

	err = writeFieldsYml(module, fileset, d)
	if err != nil {
		return err
	}

	return nil
}
