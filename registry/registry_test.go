package registry

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	tmplDir string
	regDir  string
	destDir string
	reg     *Registry
)

func TestMain(m *testing.M) {
	var err error
	destDir, err = ioutil.TempDir(os.TempDir(), "tmpl-dest")
	if err != nil {
		panic(err)
	}

	// Set up template
	setupTemplate()

	// Set up registry
	setupRegistry()

	status := m.Run()

	if err = os.RemoveAll(destDir); err != nil {
		fmt.Printf("could not clean up temp directory %s\n", destDir)
	}
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
	if _, err := reg.SaveTemplate("test", tmplDir); err != nil {
		t.Log("could not save template")
		t.FailNow()
	}
}

func testGet(t *testing.T) {
	_, err := reg.GetTemplate("test")
	if err != nil {
		t.Logf("could not get template")
		t.FailNow()
	}
}

func testGetFail(t *testing.T) {
	_, err := reg.GetTemplate("fail")
	if !IsErrTemplateNotFound(err) {
		t.Logf("template should not exist")
		t.FailNow()
	}
}

func setupTemplate() {
	var err error
	tmplDir, err = ioutil.TempDir(os.TempDir(), "tmpl-setup")
	if err != nil {
		panic(err)
	}

	// Template config
	fi, err := os.Create(filepath.Join(tmplDir, "template.toml"))
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
	regDir, err = ioutil.TempDir(os.TempDir(), "tmpl-reg")
	if err != nil {
		panic(err)
	}

	reg, err = Open(regDir)
	if err != nil {
		panic(err)
	}
}
