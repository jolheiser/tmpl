package registry

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

var (
	tmplDir string
	regDir  string
	reg     *Registry
)

func TestMain(m *testing.M) {
	var err error

	// Set up template
	setupTemplate()

	// Set up registry
	setupRegistry()

	status := m.Run()

	if err = os.RemoveAll(tmplDir); err != nil {
		fmt.Printf("could not clean up temp directory %s\n", tmplDir)
	}
	if err = os.RemoveAll(regDir); err != nil {
		fmt.Printf("could not clean up temp directory %s\n", regDir)
	}

	os.Exit(status)
}

func TestTemplate(t *testing.T) {
	t.Run("save", testSave)
	t.Run("get", testGet)
	t.Run("get-fail", testGetFail)
	t.Run("execute", testExecute)
}

func testSave(t *testing.T) {
	assert := is.New(t)
	_, err := reg.SaveTemplate("test", tmplDir)
	assert.NoErr(err) // Should save template
}

func testGet(t *testing.T) {
	assert := is.New(t)
	_, err := reg.GetTemplate("test")
	assert.NoErr(err) // Should get template
}

func testGetFail(t *testing.T) {
	assert := is.New(t)
	_, err := reg.GetTemplate("fail")
	if !errors.As(err, &ErrTemplateNotFound{}) {
		assert.Fail() // Template should not exist
	}
}

func setupTemplate() {
	var err error
	tmplDir, err = os.MkdirTemp(os.TempDir(), "tmpl-setup")
	if err != nil {
		panic(err)
	}

	// Template config
	fi, err := os.Create(filepath.Join(tmplDir, "tmpl.yaml"))
	if err != nil {
		panic(err)
	}
	_, err = fi.WriteString(tmplTemplate)
	if err != nil {
		panic(err)
	}
	if err := fi.Close(); err != nil {
		panic(err)
	}

	// Template directories
	pkgPath := filepath.Join(tmplDir, "template", "{{upper package}}")
	if err := os.MkdirAll(pkgPath, os.ModePerm); err != nil {
		panic(err)
	}
	// .tmplkeep file
	fi, err = os.Create(filepath.Join(pkgPath, ".tmplkeep"))
	if err != nil {
		panic(err)
	}
	if err := fi.Close(); err != nil {
		panic(err)
	}

	// Template file
	fi, err = os.Create(filepath.Join(tmplDir, "template", "TEST"))
	if err != nil {
		panic(err)
	}
	_, err = fi.WriteString(tmplContents)
	if err != nil {
		panic(err)
	}
	if err := fi.Close(); err != nil {
		panic(err)
	}
}

func setupRegistry() {
	var err error
	regDir, err = os.MkdirTemp(os.TempDir(), "tmpl-reg")
	if err != nil {
		panic(err)
	}

	reg, err = Open(regDir)
	if err != nil {
		panic(err)
	}
}
