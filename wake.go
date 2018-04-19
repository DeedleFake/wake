package main

import (
	"flag"
	"fmt"
	"os"
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

	file, err := openWakefile(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open Wakefile: %v", err)
		os.Exit(1)
	}
	defer file.Close()
}
