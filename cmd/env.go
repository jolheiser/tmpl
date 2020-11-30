package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
	"go.jolheiser.com/beaver/color"
)

var Env = &cli.Command{
	Name:        "env",
	Usage:       "Show tmpl environment variables",
	Description: "Show tmpl environment variables and their configuration",
	Action:      runEnv,
}

func runEnv(_ *cli.Context) error {

	// Source
	beaver.Infof("TMPL_SOURCE: %s", getEnv("TMPL_SOURCE"))

	// Registry Path
	beaver.Infof("TMPL_REGISTRY: %s", getEnv("TMPL_REGISTRY"))

	// Branch
	beaver.Infof("TMPL_BRANCH: %s", getEnv("TMPL_BRANCH"))

	return nil
}

func getEnv(key string) string {
	return color.FgHiBlue.Format(os.Getenv(key))
}
