// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

// Command vier-cli is four in a row in the terminal.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	options, err := parseOptions(os.Args[1:], os.Stderr)
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		os.Exit(2)
	}
	if options.Version {
		fmt.Println(versionLine())
		return
	}
	path, err := statePath()
	if err != nil {
		fmt.Printf("Games won't be saved: %v\n", err)
	}
	memory, err := LoadMemory(path)
	if err != nil {
		fmt.Printf("The saved game in %s could not be read, so this is a fresh start.\n", path)
	}
	NewSession(os.Stdout, path, options.apply(memory), options.Seed).Serve(os.Stdin)
}
