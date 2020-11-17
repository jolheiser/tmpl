package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var Test = &cli.Command{
	Name:        "test",
	Usage:       "Test if a directory is a valid template",
	Description: "Test whether the current directory is valid for use with tmpl",
	Action:      runTest,
}

func runTest(_ *cli.Context) error {
	var errs []string
	if _, err := os.Lstat("template.toml"); err != nil {
		errs = append(errs, "could not find template.toml")
	}

	fi, err := os.Lstat("template")
	if err != nil {
		errs = append(errs, "no template directory found")
	}
	if err == nil && !fi.IsDir() {
		errs = append(errs, "template path is a file, not a directory")
	}

	if len(errs) > 0 {
		for _, err := range errs {
			beaver.Error(err)
		}
		return nil
	}
	beaver.Info("this is a valid tmpl template")
	return nil
}
