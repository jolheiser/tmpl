package cmd

import (
	"os"

	"go.jolheiser.com/tmpl/env"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var Env = &cli.Command{
	Name:        "env",
	Usage:       "Show tmpl environment variables",
	Description: "Show tmpl environment variables and their configuration",
	Action:      runEnv,
	Subcommands: []*cli.Command{
		{
			Name:        "set",
			Aliases:     []string{"add"},
			Usage:       "Set a tmpl environment variable (stored plaintext)",
			Description: "Set an environment variable that will be loaded specifically when running tmpl (stored plaintext)",
			ArgsUsage:   "[key] [value]",
			Action:      runEnvSet,
		},
		{
			Name:        "unset",
			Aliases:     []string{"delete"},
			Usage:       "Unsets a tmpl environment variable",
			Description: "Unsets an environment variable previously set for tmpl",
			ArgsUsage:   "[key]",
			Action:      runEnvUnset,
		},
	},
}

func runEnv(_ *cli.Context) error {
	// Source
	log.Info().Str("TMPL_SOURCE", os.Getenv("TMPL_SOURCE")).Msg("system")

	// Registry Path
	log.Info().Str("TMPL_REGISTRY", os.Getenv("TMPL_REGISTRY")).Msg("system")

	// Branch
	log.Info().Str("TMPL_BRANCH", os.Getenv("TMPL_BRANCH")).Msg("system")

	// Custom
	e, err := env.Load(registryFlag)
	if err != nil {
		return err
	}
	for key, val := range e {
		log.Info().Str(key, val).Msg("env.json")
	}

	return nil
}

func runEnvSet(ctx *cli.Context) error {
	if ctx.NArg() < 2 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}
	e, err := env.Load(registryFlag)
	if err != nil {
		return err
	}
	key, val := ctx.Args().Get(0), ctx.Args().Get(1)
	e[key] = val
	if err := env.Save(registryFlag, e); err != nil {
		return err
	}
	log.Info().Str(key, val).Msg("Successfully saved tmpl environment variable!")
	return nil
}

func runEnvUnset(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}
	e, err := env.Load(registryFlag)
	if err != nil {
		return err
	}
	key := ctx.Args().First()
	val := e[key]
	delete(e, key)
	if err := env.Save(registryFlag, e); err != nil {
		return err
	}
	log.Info().Str(key, val).Msg("Successfully unset tmpl environment variable!")
	return nil
}
