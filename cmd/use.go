package cmd

import (
	"errors"

	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var Use = &cli.Command{
	Name:        "use",
	Usage:       "Use a template",
	Description: "Use (execute) a template from the registry",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "defaults",
			Usage: "Use template defaults",
		},
	},
	Action: runUse,
}

func runUse(ctx *cli.Context) error {
	if ctx.NArg() < 2 {
		return errors.New("<name> <dest>")
	}

	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	tmpl, err := reg.GetTemplate(ctx.Args().First())
	if err != nil {
		return err
	}

	if err := tmpl.Execute(ctx.Args().Get(1), ctx.Bool("defaults")); err != nil {
		return err
	}

	beaver.Infof("Successfully executed %s", tmpl.Name)
	return nil
}
