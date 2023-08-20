package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

var (
	tmplContents = `{{title name}} (@{{username}}) {{if .bool}}{{.year}}{{end}} {{org}}`
	tmplTemplate = `
prompts:
  - id: name
    default: john olheiser
  - id: year
    default: ${TMPL_TEST} # 2020
  - id: package
    default: pkg
  - id: bool
    default: true
  - id: username
    default: username
  - id: org
    default: ${TMPL_PROMPT_USERNAME}/org
`
	tmplGold    = "John Olheiser (@jolheiser) 2020 jolheiser/org"
	tmplNewGold = "DO NOT OVERWRITE!"
)

func testExecute(t *testing.T) {
	assert := is.New(t)
	destDir := t.TempDir()

	// Set environment variable
	err := os.Setenv("TMPL_TEST", "2020")
	assert.NoErr(err) // Should set TMPL_TEST env

	err = os.Setenv("TMPL_VAR_USERNAME", "jolheiser")
	assert.NoErr(err) // Should set TMPL_VAR_USERNAME env

	// Get template
	tmpl, err := reg.GetTemplate("test")
	assert.NoErr(err) // Should get template

	// Execute template
	err = tmpl.Execute(destDir, true, true)
	assert.NoErr(err) // Should execute template

	// Check contents of file
	testPath := filepath.Join(destDir, "TEST")
	contents, err := os.ReadFile(testPath)
	assert.NoErr(err)                        // Should be able to read TEST file
	assert.Equal(string(contents), tmplGold) // Template should match golden file

	// Check if directory was created
	pkgPath := filepath.Join(destDir, "PKG")
	_, err = os.Lstat(pkgPath)
	assert.NoErr(err) // PKG directory should exist

	// Check for .tmplkeep
	tmplKeep := filepath.Join(pkgPath, ".tmplkeep")
	_, err = os.Lstat(tmplKeep)
	assert.True(err != nil) // .tmplkeep file should NOT be retained

	// Change file to test non-overwrite
	err = os.WriteFile(testPath, []byte(tmplNewGold), os.ModePerm)
	assert.NoErr(err) // Writing file should succeed

	err = tmpl.Execute(destDir, true, false)
	assert.NoErr(err) // Should execute template

	contents, err = os.ReadFile(testPath)
	assert.NoErr(err)                           // Should be able to read file
	assert.Equal(string(contents), tmplNewGold) // Template should match new golden file
}
