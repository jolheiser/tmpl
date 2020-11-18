package cmd

import (
	"errors"

	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var Update = &cli.Command{
	Name:        "update",
	Usage:       "Update a template",
	Description: "Update a template in the registry from the original source",
	Action:      runUpdate,
}

func runUpdate(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return errors.New("<name>")
	}

	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	tmpl, err := reg.GetTemplate(ctx.Args().First())
	if err != nil {
		return err
	}

	if err := reg.RemoveTemplate(tmpl.Name); err != nil {
		return err
	}

	if tmpl.Path != "" {
		_, err = reg.SaveTemplate(tmpl.Name, tmpl.Path)
	} else {
		_, err = reg.DownloadTemplate(tmpl.Name, tmpl.Repository, tmpl.Branch)
	}
	if err != nil {
		return err
	}

	beaver.Infof("Successfully updated %s", tmpl.Name)
	return nil
}
