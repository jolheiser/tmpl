package cmd

import (
	"errors"
	"path/filepath"

	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var Save = &cli.Command{
	Name:        "save",
	Usage:       "Save a local template",
	Description: "Save a local template to the registry",
	Action:      runSave,
}

func runSave(ctx *cli.Context) error {
	if ctx.NArg() < 2 {
		return errors.New("<path> <name>")
	}

	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	localPath := ctx.Args().First()
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return err
	}

	t, err := reg.SaveTemplate(ctx.Args().Get(1), localPath)
	if err != nil {
		return err
	}

	beaver.Infof("Added new template %s", t.Name)
	return nil
}
