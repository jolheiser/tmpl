package cmd

import (
	"os"
	"path/filepath"

	"go.jolheiser.com/tmpl/cmd/flags"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var (
	Version    = "develop"
	defaultDir string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		beaver.Error("could not locate user's home directory, tmpl will use temp dir for registry")
		return
	}
	defaultDir = filepath.Join(home, ".tmpl")
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "tmpl"
	app.Usage = "Template automation"
	app.Description = "Template automation"
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "registry",
			Aliases:     []string{"r"},
			Usage:       "Registry directory of tmpl",
			Value:       defaultDir,
			DefaultText: "~/.tmpl",
			Destination: &flags.Registry,
			EnvVars:     []string{"TMPL_REGISTRY"},
		},
		&cli.StringFlag{
			Name:        "source",
			Aliases:     []string{"s"},
			Usage:       "Short-name source to use",
			Destination: &flags.Source,
			EnvVars:     []string{"TMPL_SOURCE"},
		},
	}

	app.Commands = []*cli.Command{
		Download,
		Init,
		List,
		Remove,
		Save,
		Source,
		Test,
		Update,
		Use,
	}

	return app
}
