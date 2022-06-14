package cmd

import (
	"go.jolheiser.com/tmpl/registry"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var Remove = &cli.Command{
	Name:        "remove",
	Usage:       "Remove a template",
	Description: "Remove a template from the registry",
	ArgsUsage:   "[name]",
	Action:      runRemove,
}

func runRemove(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}

	reg, err := registry.Open(registryFlag)
	if err != nil {
		return err
	}

	if err := reg.RemoveTemplate(ctx.Args().First()); err != nil {
		return err
	}

	log.Info().Msgf("Successfully removed %q", ctx.Args().First())
	return nil
}
