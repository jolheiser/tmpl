package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"go.jolheiser.com/tmpl/registry"

	"github.com/urfave/cli/v2"
)

var List = &cli.Command{
	Name:        "list",
	Usage:       "List templates in the registry",
	Description: "List all usable templates currently downloaded in the registry",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "json",
			Usage: "JSON format",
		},
	},
	Action: runList,
}

func runList(ctx *cli.Context) error {
	reg, err := registry.Open(registryFlag)
	if err != nil {
		return err
	}

	if ctx.Bool("json") {
		return json.NewEncoder(os.Stdout).Encode(reg.Templates)
	}

	wr := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	if _, err := fmt.Fprintf(wr, "NAME\tURL\tLOCAL\tLAST UPDATED\n"); err != nil {
		return err
	}
	for _, t := range reg.Templates {
		u := fmt.Sprintf("%s @%s", t.Repository, t.Branch)
		var local bool
		if t.Path != "" {
			u = t.Path
			local = true
		}
		if _, err := fmt.Fprintf(wr, "%s\t%s\t%t\t%s\n", t.Name, u, local, t.LastUpdate.Format("01/02/2006")); err != nil {
			return err
		}
	}
	return wr.Flush()
}
