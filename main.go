package main

import (
	"os"

	"go.jolheiser.com/tmpl/cmd"

	"go.jolheiser.com/beaver"
	"go.jolheiser.com/beaver/color"
)

func main() {
	app := cmd.NewApp()
	color.Fatal = color.Error // Easier to read, doesn't need to stand out as much in a CLI
	if err := app.Run(os.Args); err != nil {
		beaver.Fatal(err)
	}
}
