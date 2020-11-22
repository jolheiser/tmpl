package registry

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	tmplContents = `{{title name}} {{.year}}`
	tmplTemplate = `name = "john olheiser"

[year]
default = 2020

[package]
default = "pkg"`
	tmplGold    = "John Olheiser 2020"
	tmplNewGold = "DO NOT OVERWRITE!"
)

func testExecute(t *testing.T) {
	// Get template
	tmpl, err := reg.GetTemplate("test")
	if err != nil {
		t.Logf("could not get template")
		t.FailNow()
	}

	// Execute template
	if err := tmpl.Execute(destDir, true, true); err != nil {
		t.Logf("could not execute template: %v\n", err)
		t.FailNow()
	}

	// Check contents of file
	testPath := filepath.Join(destDir, "TEST")
	contents, err := ioutil.ReadFile(testPath)
	if err != nil {
		t.Logf("could not read file: %v\n", err)
		t.FailNow()
	}

	if string(contents) != tmplGold {
		t.Logf("contents did not match:\n\tExpected: %s\n\tGot: %s", tmplGold, string(contents))
		t.FailNow()
	}

	// Check if directory was created
	pkgPath := filepath.Join(destDir, "PKG")
	if _, err := os.Lstat(pkgPath); err != nil {
		t.Logf("expected a directory at %s: %v\n", pkgPath, err)
		t.FailNow()
	}

	// Change file to test non-overwrite
	if err := ioutil.WriteFile(testPath, []byte(tmplNewGold), os.ModePerm); err != nil {
		t.Logf("could not write file: %v\n", err)
		t.FailNow()
	}

	if err := tmpl.Execute(destDir, true, false); err != nil {
		t.Logf("could not execute template: %v\n", err)
		t.FailNow()
	}

	contents, err = ioutil.ReadFile(testPath)
	if err != nil {
		t.Logf("could not read file: %v\n", err)
		t.FailNow()
	}

	if string(contents) != tmplNewGold {
		t.Logf("contents did not match:\n\tExpected: %s\n\tGot: %s", tmplNewGold, string(contents))
		t.FailNow()
	}
}
