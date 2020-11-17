// +build docs

package main

import (
	"os"
	"strings"

	"go.jolheiser.com/tmpl/cmd"
)

func main() {
	app := cmd.NewApp()

	fi, err := os.Create("DOCS.md")
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	md, err := app.ToMarkdown()
	if err != nil {
		panic(err)
	}

	// CLI ToMarkdown issue related to man-pages
	md = md[strings.Index(md, "#"):]

	if _, err := fi.WriteString(md); err != nil {
		panic(err)
	}
}
