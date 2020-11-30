package cmd

import (
	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var Update = &cli.Command{
	Name:        "update",
	Usage:       "Update a template",
	Description: "Update a template in the registry from the original source",
	ArgsUsage:   "[name]",
	Action:      runUpdate,
}

func runUpdate(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}

	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	tmpl, err := reg.GetTemplate(ctx.Args().First())
	if err != nil {
		return err
	}

	if err := reg.UpdateTemplate(tmpl.Name); err != nil {
		return err
	}

	beaver.Infof("Successfully updated %s", tmpl.Name)
	return nil
}
