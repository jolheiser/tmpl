package cmd

import (
	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
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

	log.Info().Msgf("Successfully updated %q", tmpl.Name)
	return nil
}
