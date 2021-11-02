package cmd

import (
	"os"

	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
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
			log.Info().Msgf("Restoring %q...", tmpl.Name)
			if err := reg.UpdateTemplate(tmpl.Name); err != nil {
				return err
			}
			num++
		}
	}

	log.Info().Int("count", num).Msgf("Restored templates.")
	return nil
}
