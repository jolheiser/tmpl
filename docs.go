//go:build docs
// +build docs

package main

import (
	"os"
	"regexp"
	"strings"

	"go.jolheiser.com/tmpl/cmd"
)

//go:generate go run docs.go
func main() {
	app := cmd.NewApp()

	fi, err := os.Create("CLI.md")
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

	// CLI is using real default instead of DefaultText
	md = regexp.MustCompile(`[\/\\:\w]+\.tmpl`).ReplaceAllString(md, "~/.tmpl")

	if _, err := fi.WriteString(md); err != nil {
		panic(err)
	}
}
