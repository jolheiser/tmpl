package cmd

import (
	"errors"
	"path/filepath"
	"strings"

	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var Save = &cli.Command{
	Name:        "save",
	Usage:       "Save a local template",
	Description: "Save a local template to the registry",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Branch to clone",
			Value:   "main",
			EnvVars: []string{"TMPL_BRANCH"},
		},
	},
	Action: runSave,
}

func runSave(ctx *cli.Context) error {
	if ctx.NArg() < 2 {
		return errors.New("<path> <name>")
	}

	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	// Did the user give us the root path, or the .git directory?
	localPath := ctx.Args().First()
	if !strings.HasSuffix(localPath, ".git") {
		localPath = filepath.Join(localPath, ".git")
	}
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return err
	}

	t, err := reg.AddTemplate(ctx.Args().Get(1), localPath, ctx.String("branch"))
	if err != nil {
		return err
	}

	beaver.Infof("Added new template %s", t.Name)
	return nil
}
