package cmd

import (
	"path/filepath"

	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var Save = &cli.Command{
	Name:        "save",
	Usage:       "Save a local template",
	Description: "Save a local template to the registry",
	ArgsUsage:   "[path] [name]",
	Action:      runSave,
}

func runSave(ctx *cli.Context) error {
	if ctx.NArg() < 2 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
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

	log.Info().Msgf("Added new template %q", t.Name)
	return nil
}
