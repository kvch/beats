package fileset

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	modulePath = "../"
	beatsPath  = "."
)

func copyTemplatesToDest(templatesPath, name, filesetPath, module, fileset string) error {
	template := path.Join(templatesPath, name)
	c, err := ioutil.ReadFile(template)
	if err != nil {
		return err
	}

	c = bytes.Replace(c, []byte("{module}"), []byte(module), -1)
	c = bytes.Replace(c, []byte("{fileset}"), []byte(fileset), -1)

	dest := path.Join(filesetPath, name)
	err = ioutil.WriteFile(dest, c, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot copy template: %v", err)
	}
	return nil
}

func generateModule(module, fileset string) error {
	p := path.Join(modulePath, "module", module)
	if _, err := os.Stat(p); os.IsExist(err) {
		return fmt.Errorf("module already exists: %s at %s", module, p)
	}

	d := path.Join(p, "_meta", "kibana", "default")
	err := os.MkdirAll(d, 0750)
	if err != nil {
		return err
	}

	templatesPath := path.Join(beatsPath, "module")
	filesToCopy := []string{path.Join("fields.yml"), path.Join("docs.asciidoc")}
	for _, f := range filesToCopy {
		err := copyTemplatesToDest(templatesPath, f, p, module, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func generateFileset(module, fileset string) error {
	filesetPath := path.Join(modulePath, "module", module, fileset)
	if _, err := os.Stat(filesetPath); os.IsExist(err) {
		return fmt.Errorf("fileset already exists: %s", fileset)
	}

	dirsToCreate := []string{"", "_meta", "test", "config", "ingest"}
	for _, d := range dirsToCreate {
		p := path.Join(filesetPath, d)
		err := os.Mkdir(p, 0750)
		if err != nil {
			return err
		}
	}

	templatesPath := path.Join(beatsPath, "module", "fileset")
	filesToCopy := []string{path.Join("config", "config.yml"), path.Join("ingest", "pipeline.json"), "manifest.yml"}
	for _, f := range filesToCopy {
		err := copyTemplatesToDest(templatesPath, f, filesetPath, module, fileset)
		if err != nil {
			return err
		}
	}
	return nil
}

// Generate creates required folders and configuration files. Names of fileset and modules are
// substituted in these created files.
func Generate(module, fileset string) error {
	err := generateModule(module, fileset)
	if err != nil {
		return err
	}

	err = generateFileset(module, fileset)
	if err != nil {
		return err
	}

	return nil
}
