// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

// Command vier-cli is four in a row in the terminal.
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	path, err := statePath()
	if err != nil {
		fmt.Printf("Games won't be saved: %v\n", err)
	}
	memory, err := LoadMemory(path)
	if err != nil {
		fmt.Printf("The saved game in %s could not be read, so this is a fresh start.\n", path)
	}
	NewSession(os.Stdout, path, memory, uint64(time.Now().UnixNano())).Serve(os.Stdin)
}
