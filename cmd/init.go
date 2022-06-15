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
	if _, err := os.Lstat("tmpl.yaml"); !os.IsNotExist(err) {
		if err != nil {
			return err
		}
		return errors.New("tmpl.yaml already detected, aborting initialization")
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

	fi, err := os.Create("tmpl.yaml")
	if err != nil {
		return err
	}
	if _, err := fi.WriteString(initConfig); err != nil {
		return err
	}
	if err := os.Mkdir("template", os.ModePerm); err != nil {
		return err
	}
	log.Info().Msg("Template initialized!")
	return fi.Close()
}

var initConfig = `# tmpl.yaml
# Write any template args here to prompt the user for, giving any defaults/options as applicable

prompts:
  - id: name                              # The unique ID for the prompt
    label: Project Name                   # (Optional) Prompt message/label, defaults to id
    help: The name to use in the project  # (Optional) Help message for the prompt
    default: tmpl                         # (Optional) Prompt default
    options:                              # (Optional) Set of options the user can choose from
      - coolproject123
      - ${USER}'s cool project
`
