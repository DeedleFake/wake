package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/all"
)

func openWakefile(name string) (*os.File, error) {
	for {
		file, err := os.Open(name)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.Chdir("..")
				if err != nil {
					return nil, err
				}

				continue
			}

			return nil, err
		}
		return file, nil
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] [target]\n\n", os.Args[0])

		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	filename := flag.String("f", "Wakefile", "Name of the file to load targets from.")
	flag.Parse()

	target := flag.Arg(0)
	if target == "" {
		target = "default"
	}

	file, err := openWakefile(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open Wakefile: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	c, err := wdte.Parse(file, std.Import)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse %q: %v", *filename, err)
		os.Exit(1)
	}

	scope, r := c.Collect(std.F())
	if err, ok := r.(error); ok {
		fmt.Fprintf(os.Stderr, "Wakefile scope collection failed: %v", err)
		os.Exit(1)
	}

	// TODO: Implement target patterns.
	rule := scope.Get(wdte.ID(target))
	if rule == nil {
		fmt.Fprintf(os.Stderr, "No rule for %q", target)
		os.Exit(1)
	}

	r = rule.Call(std.F())
	if err, ok := r.(error); ok {
		fmt.Fprintf(os.Stderr, "Error during rule execution: %v", err)
		os.Exit(1)
	}
}
