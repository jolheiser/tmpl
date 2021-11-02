package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var Env = &cli.Command{
	Name:        "env",
	Usage:       "Show tmpl environment variables",
	Description: "Show tmpl environment variables and their configuration",
	Action:      runEnv,
}

func runEnv(_ *cli.Context) error {

	// Source
	log.Info().Str("TMPL_SOURCE", os.Getenv("TMPL_SOURCE")).Msg("")

	// Registry Path
	log.Info().Str("TMPL_REGISTRY", os.Getenv("TMPL_REGISTRY")).Msg("")

	// Branch
	log.Info().Str("TMPL_BRANCH", os.Getenv("TMPL_BRANCH")).Msg("")

	return nil
}


