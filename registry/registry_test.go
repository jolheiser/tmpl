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

	tmplContents = `{{title name}} {{year}}`
	tmplTemplate = `name = "john olheiser"
year = 2020`
	tmplGold = "John Olheiser 2020"
)

func TestMain(m *testing.M) {
	var err error
	destDir, err = ioutil.TempDir(os.TempDir(), "tmpl")
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

func testExecute(t *testing.T) {
	tmpl, err := reg.GetTemplate("test")
	if err != nil {
		t.Logf("could not get template")
		t.FailNow()
	}

	if err := tmpl.Execute(destDir, true); err != nil {
		t.Logf("could not execute template: %v\n", err)
		t.FailNow()
	}

	contents, err := ioutil.ReadFile(filepath.Join(destDir, "TEST"))
	if err != nil {
		t.Logf("could not read file: %v\n", err)
		t.FailNow()
	}

	if string(contents) != tmplGold {
		t.Logf("contents did not match:\n\tExpected: %s\n\tGot: %s", tmplGold, string(contents))
		t.FailNow()
	}
}

func setupTemplate() {
	var err error
	tmplDir, err = ioutil.TempDir(os.TempDir(), "tmpl")
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

	// Template file
	if err := os.Mkdir(filepath.Join(tmplDir, "template"), os.ModePerm); err != nil {
		panic(err)
	}
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
	regDir, err = ioutil.TempDir(os.TempDir(), "tmpl")
	if err != nil {
		panic(err)
	}

	reg, err = Open(regDir)
	if err != nil {
		panic(err)
	}
}
