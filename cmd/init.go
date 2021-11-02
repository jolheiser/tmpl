package cmd

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var Init = &cli.Command{
	Name:        "init",
	Usage:       "Initialize a template",
	Description: "Initializes a template structure for creating a new tmpl template",
	Action:      runInit,
}

func runInit(_ *cli.Context) error {
	if _, err := os.Lstat("template.toml"); !os.IsNotExist(err) {
		if err != nil {
			return err
		}
		return errors.New("template.toml already detected, aborting initialization")
	}
	if fi, err := os.Lstat("template"); !os.IsNotExist(err) {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			return errors.New("template file found instead of directory, aborting initialization")
		}
		return errors.New("template directory already detected, aborting initialization")
	}

	fi, err := os.Create("template.toml")
	if err != nil {
		return err
	}
	if _, err := fi.WriteString(comments); err != nil {
		return err
	}
	if err := os.Mkdir("template", os.ModePerm); err != nil {
		return err
	}
	log.Info().Msg("Template initialized!")
	return fi.Close()
}

var comments = `# template.toml
# Write any template args here to prompt the user for, giving any defaults/options as applicable

[name]
prompt = "Project Name"
help = "The name to use in the project"
default = "tmpl"
`
