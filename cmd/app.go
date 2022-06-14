package cmd

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	Version    = "develop"
	defaultDir string

	registryFlag string
	sourceFlag   string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Error().Msg("could not locate user's home directory, tmpl will use temp dir for registry")
		defaultDir = filepath.Join(os.TempDir(), ".tmpl")
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
			Destination: &registryFlag,
			EnvVars:     []string{"TMPL_REGISTRY"},
		},
		&cli.StringFlag{
			Name:        "source",
			Aliases:     []string{"s"},
			Usage:       "Short-name source to use",
			Destination: &sourceFlag,
			EnvVars:     []string{"TMPL_SOURCE"},
		},
	}

	app.Commands = []*cli.Command{
		Download,
		Env,
		Init,
		List,
		Remove,
		Restore,
		Save,
		Source,
		Test,
		Update,
		Use,
	}

	return app
}
