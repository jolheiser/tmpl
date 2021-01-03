package registry

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	tmplContents = `{{title name}} (@{{username}}) {{if .bool}}{{.year}}{{end}}`
	tmplTemplate = `
name = "john olheiser"

[year]
default = ${TMPL_TEST} # 2020

[package]
default = "pkg"

[bool]
default = true

[username]
default = "username"
`
	tmplGold    = "John Olheiser (@jolheiser) 2020"
	tmplNewGold = "DO NOT OVERWRITE!"
)

func testExecute(t *testing.T) {
	// Set environment variable
	if err := os.Setenv("TMPL_TEST", "2020"); err != nil {
		t.Logf("could not set environment: %v", err)
		t.FailNow()
	}
	if err := os.Setenv("TMPL_VAR_USERNAME", "jolheiser"); err != nil {
		t.Logf("could not set environment: %v", err)
		t.FailNow()
	}

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

	// Check for .tmplkeep
	tmplKeep := filepath.Join(pkgPath, ".tmplkeep")
	if _, err := os.Lstat(tmplKeep); err == nil {
		t.Logf(".tmplkeep files should NOT be retained upon execution: %s\n", tmplKeep)
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
