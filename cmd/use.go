package cmd

import (
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
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Overwrite existing files",
		},
	},
	ArgsUsage: "[name] [destination (default: \".\")]",
	Action:    runUse,
}

func runUse(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}

	dest := "."
	if ctx.NArg() >= 2 {
		dest = ctx.Args().Get(1)
	}

	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	tmpl, err := reg.GetTemplate(ctx.Args().First())
	if err != nil {
		return err
	}

	if err := tmpl.Execute(dest, ctx.Bool("defaults"), ctx.Bool("force")); err != nil {
		return err
	}

	beaver.Infof("Successfully executed %s", tmpl.Name)
	return nil
}
