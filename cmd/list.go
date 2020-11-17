package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/urfave/cli/v2"
)

var List = &cli.Command{
	Name:        "list",
	Usage:       "List templates in the registry",
	Description: "List all usable templates currently downloaded in the registry",
	Action:      runList,
}

func runList(_ *cli.Context) error {
	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	wr := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	for _, t := range reg.Templates {
		if _, err := fmt.Fprintf(wr, "%s\t%s@%s\t%s\n", t.Name, t.Repository, t.Branch, t.Created); err != nil {
			return err
		}
	}
	return wr.Flush()
}
