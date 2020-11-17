package registry

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mholt/archiver/v3"
	"github.com/pelletier/go-toml"
)

// Template is a tmpl project
type Template struct {
	reg        *Registry `toml:"-"`
	Name       string    `toml:"name"`
	Repository string    `toml:"repository"`
	Branch     string    `toml:"branch"`
	Created    time.Time `toml:"created"`
}

// ArchiveName is the name given to the archive for this Template
func (t *Template) ArchiveName() string {
	return fmt.Sprintf("%s.tar.gz", t.Name)
}

// ArchivePath is the full path to the archive for this Template within the Registry
func (t *Template) ArchivePath() string {
	return filepath.Join(t.reg.dir, t.ArchiveName())
}

// Execute runs the Template and copies to dest
func (t *Template) Execute(dest string) error {
	tmp, err := ioutil.TempDir(os.TempDir(), "tmpl")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	if err := archiver.Unarchive(t.ArchivePath(), tmp); err != nil {
		return err
	}

	vars, err := prompt(tmp)
	if err != nil {
		return err
	}

	base := filepath.Join(tmp, "template")
	return filepath.Walk(base, func(walkPath string, walkInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if walkInfo.IsDir() {
			return nil
		}

		contents, err := ioutil.ReadFile(walkPath)
		if err != nil {
			return err
		}

		tmpl, err := template.New("tmpl").Parse(string(contents))
		if err != nil {
			return err
		}

		newDest := strings.TrimPrefix(walkPath, base+"/")
		newDest = filepath.Join(dest, newDest)

		if err := os.MkdirAll(filepath.Dir(newDest), os.ModePerm); err != nil {
			return err
		}

		oldFi, err := os.Lstat(walkPath)
		if err != nil {
			return err
		}
		newFi, err := os.OpenFile(newDest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, oldFi.Mode())
		if err != nil {
			return err
		}

		if err := tmpl.Execute(newFi, vars); err != nil {
			return err
		}

		return newFi.Close()
	})
}

func prompt(dir string) (map[string]interface{}, error) {
	templatePath := filepath.Join(dir, "template.toml")
	if _, err := os.Lstat(templatePath); err != nil {
		return nil, err
	}

	tree, err := toml.LoadFile(templatePath)
	if err != nil {
		return nil, err
	}

	// Sort the map keys so they are consistent
	vars := tree.ToMap()
	sorted := make([]string, 0, len(vars))
	for k := range vars {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	for _, k := range sorted {
		v := vars[k]
		var p survey.Prompt
		switch t := v.(type) {
		case []string:
			p = &survey.Select{
				Message: k,
				Options: t,
			}
		default:
			p = &survey.Input{
				Message: k,
				Default: fmt.Sprintf("%v", t),
			}
		}
		q := []*survey.Question{
			{
				Name:     "response",
				Prompt:   p,
				Validate: survey.Required,
			},
		}
		a := struct {
			Response string
		}{}
		if err := survey.Ask(q, &a); err != nil {
			return nil, err
		}
		vars[k] = a.Response
	}

	return vars, nil
}
