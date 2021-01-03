package registry

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/mholt/archiver/v3"
)

// Template is a tmpl project
type Template struct {
	reg        *Registry `toml:"-"`
	Name       string    `toml:"name"`
	Path       string    `toml:"path"`
	Repository string    `toml:"repository"`
	Branch     string    `toml:"branch"`
	LastUpdate time.Time `toml:"last_update"`
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
func (t *Template) Execute(dest string, defaults, overwrite bool) error {
	tmp, err := ioutil.TempDir(os.TempDir(), "tmpl")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	if err := archiver.Unarchive(t.ArchivePath(), tmp); err != nil {
		return err
	}

	prompts, err := prompt(tmp, defaults)
	if err != nil {
		return err
	}

	funcs := mergeMaps(funcMap, prompts.ToFuncMap())

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

		newDest := strings.TrimPrefix(walkPath, base+string(filepath.Separator))
		newDest = filepath.Join(dest, newDest)

		tmplDest, err := template.New("dest").Funcs(funcs).Parse(newDest)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tmplDest.Execute(&buf, prompts.ToMap()); err != nil {
			return err
		}
		newDest = buf.String()

		if err := os.MkdirAll(filepath.Dir(newDest), os.ModePerm); err != nil {
			return err
		}

		// Skip .tmplkeep files, after creating the directory structure
		if strings.EqualFold(walkInfo.Name(), ".tmplkeep") {
			return nil
		}

		oldFi, err := os.Lstat(walkPath)
		if err != nil {
			return err
		}

		// Check if new file exists. If it does, only skip if not overwriting
		if _, err := os.Lstat(newDest); err == nil && !overwrite {
			return nil
		}

		newFi, err := os.OpenFile(newDest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, oldFi.Mode())
		if err != nil {
			return err
		}

		tmplContents, err := template.New("tmpl").Funcs(funcs).Parse(string(contents))
		if err != nil {
			return err
		}
		if err := tmplContents.Execute(newFi, prompts.ToMap()); err != nil {
			return err
		}

		return newFi.Close()
	})
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for _, mm := range maps {
		for k, v := range mm {
			m[k] = v
		}
	}
	return m
}
