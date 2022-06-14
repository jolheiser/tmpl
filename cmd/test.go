package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.jolheiser.com/tmpl/schema"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
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

	fi, err := os.Open(filepath.Join(testPath, "tmpl.yaml"))
	if err != nil {
		errs = append(errs, fmt.Sprintf("could not open tmpl.yaml: %v", err))
	}
	defer fi.Close()
	if err := schema.Lint(fi); err != nil {
		var rerr schema.ResultErrors
		if errors.As(err, &rerr) {
			for _, re := range rerr {
				errs = append(errs, fmt.Sprintf("%s: %s", re.Field(), re.Description()))
			}
		} else {
			errs = append(errs, fmt.Sprintf("could not lint tmpl.yaml: %v", err))
		}
	}

	fstat, err := os.Lstat(filepath.Join(testPath, "template"))
	if err != nil {
		errs = append(errs, "no template directory found")
	}
	if err == nil && !fstat.IsDir() {
		errs = append(errs, "template path is a file, not a directory")
	}

	if len(errs) > 0 {
		for _, err := range errs {
			log.Error().Msg(err)
		}
		return nil
	}

	log.Info().Msg("This is a valid tmpl template!")
	return nil
}
