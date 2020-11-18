package cmd

import (
	"errors"
	"fmt"
	"strings"

	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var Download = &cli.Command{
	Name:        "download",
	Usage:       "Download a template",
	Description: "Download a template and save it to the local registry",
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
	if ctx.NArg() < 2 {
		return errors.New("<repo> <name>")
	}

	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	var source *registry.Source
	if flags.Source != "" {
		for _, s := range reg.Sources {
			if strings.EqualFold(s.Name, flags.Source) {
				source = s
				break
			}
		}
		if source == nil {
			return fmt.Errorf("could not find source for %s", flags.Source)
		}
	}

	cloneURL := ctx.Args().First()
	if !strings.HasSuffix(cloneURL, ".git") {
		cloneURL += ".git"
	}
	if source != nil {
		cloneURL = source.CloneURL(cloneURL)
	}

	t, err := reg.DownloadTemplate(ctx.Args().Get(1), cloneURL, ctx.String("branch"))
	if err != nil {
		return err
	}

	beaver.Infof("Added new template %s", t.Name)
	return nil
}
