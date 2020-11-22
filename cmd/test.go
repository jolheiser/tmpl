package cmd

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var Test = &cli.Command{
	Name:        "test",
	Usage:       "Test if a directory is a valid template",
	Description: "Test whether a directory is valid for use with tmpl",
	ArgsUsage:   "[path (default: \".\")]",
	Action:      runTest,
}

func runTest(ctx *cli.Context) error {
	testPath := "."
	if ctx.NArg() > 0 {
		testPath = ctx.Args().First()
	}

	var errs []string
	if _, err := os.Lstat(filepath.Join(testPath, "template.toml")); err != nil {
		errs = append(errs, "could not find template.toml")
	}

	fi, err := os.Lstat(filepath.Join(testPath, "template"))
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
