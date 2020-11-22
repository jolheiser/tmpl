package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"go.jolheiser.com/tmpl/cmd/flags"
	"go.jolheiser.com/tmpl/registry"

	"github.com/urfave/cli/v2"
	"go.jolheiser.com/beaver"
)

var (
	Source = &cli.Command{
		Name:        "source",
		Usage:       "Commands for working with sources",
		Description: "Commands for working with sources, short-hand flags for easier downloads",
		Action:      SourceList.Action,
		Subcommands: []*cli.Command{
			SourceList,
			SourceAdd,
			SourceRemove,
		},
	}

	SourceList = &cli.Command{
		Name:        "list",
		Usage:       "List available sources",
		Description: "List all available sources in the registry",
		Action:      runSourceList,
	}

	SourceAdd = &cli.Command{
		Name:        "add",
		Usage:       "Add a source",
		Description: "Add a new source to the registry",
		ArgsUsage:   "[base URL] [name]",
		Action:      runSourceAdd,
	}

	SourceRemove = &cli.Command{
		Name:        "remove",
		Usage:       "Remove a source",
		Description: "Remove a source from the registry",
		ArgsUsage:   "[name]",
		Action:      runSourceRemove,
	}
)

func runSourceList(_ *cli.Context) error {
	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	wr := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	if _, err := fmt.Fprintf(wr, "NAME\tURL\n"); err != nil {
		return err
	}
	for _, s := range reg.Sources {
		if _, err := fmt.Fprintf(wr, "%s\t%s\n", s.Name, s.URL); err != nil {
			return err
		}
	}
	return wr.Flush()
}

func runSourceAdd(ctx *cli.Context) error {
	if ctx.NArg() < 2 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}

	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	s, err := reg.AddSource(ctx.Args().First(), ctx.Args().Get(1))
	if err != nil {
		return err
	}

	beaver.Infof("Added new source %s", s.Name)
	return nil
}

func runSourceRemove(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return cli.ShowCommandHelp(ctx, ctx.Command.Name)
	}

	reg, err := registry.Open(flags.Registry)
	if err != nil {
		return err
	}

	if err := reg.RemoveSource(ctx.Args().First()); err != nil {
		return err
	}

	beaver.Infof("Successfully removed source for %s", ctx.Args().First())
	return nil
}
