package main

import (
	"os"

	"go.jolheiser.com/tmpl/cmd"

	"go.jolheiser.com/beaver"
)

func main() {
	app := cmd.NewApp()

	if err := app.Run(os.Args); err != nil {
		beaver.Fatal(err)
	}
}
