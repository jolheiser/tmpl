package main

import (
	"os"

	"go.jolheiser.com/tmpl/cmd"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	app := cmd.NewApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
