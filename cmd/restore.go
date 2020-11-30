package cmd

import (
	"os"

	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var Restore = &cli.Command{
	Name:        "restore",
	Usage:       "Restore missing templates",
	Description: "Restore templates that are listed in the registry, but are missing archives",
	Action:      runRestore,
}

func runRestore(_ *cli.Context) error {
	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	var num int
	for _, tmpl := range reg.Templates {
		if _, err := os.Lstat(tmpl.ArchivePath()); os.IsNotExist(err) {
			beaver.Infof("Restoring %s...", tmpl.Name)
			if err := reg.UpdateTemplate(tmpl.Name); err != nil {
				return err
			}
			num++
		}
	}

	beaver.Infof("Restored %d templates.", num)
	return nil
}
