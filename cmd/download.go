package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"go.jolheiser.com/tmpl/env"
	"go.jolheiser.com/tmpl/registry"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var Download = &cli.Command{
	Name:        "download",
	Usage:       "Download a template",
	Description: "Download a template and save it to the local registry",
	ArgsUsage:   "[repository URL] <name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Branch to clone",
			Value:   "main",
			EnvVars: []string{"TMPL_BRANCH"},
		},
	},
	Action: runDownload,
}

func runDownload(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}

	reg, err := registry.Open(registryFlag)
	if err != nil {
		return err
	}

	e, err := env.Load(registryFlag)
	if err != nil {
		return err
	}
	if err := e.Set(); err != nil {
		return err
	}

	var source *registry.Source
	if sourceFlag != "" {
		for _, s := range reg.Sources {
			if strings.EqualFold(s.Name, sourceFlag) {
				source = s
				break
			}
		}
		if source == nil {
			return fmt.Errorf("could not find source for %s", sourceFlag)
		}
	}

	cloneURL := ctx.Args().First()
	if source != nil {
		cloneURL = source.CloneURL(cloneURL)
	}
	if !strings.HasSuffix(cloneURL, ".git") {
		cloneURL += ".git"
	}

	t, err := reg.DownloadTemplate(deriveName(ctx), cloneURL, ctx.String("branch"))
	if err != nil {
		return err
	}

	log.Info().Msgf("Added new template %q", t.Name)
	return nil
}

func deriveName(ctx *cli.Context) string {
	if ctx.NArg() > 1 {
		return ctx.Args().Get(1)
	}

	envBranch, envSet := os.LookupEnv("TMPL_BRANCH")
	flagBranch, flagSet := ctx.String("branch"), ctx.IsSet("branch")
	if flagSet {
		if !envSet || envBranch != flagBranch {
			return flagBranch
		}
	}

	return path.Base(ctx.Args().First())
}
