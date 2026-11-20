// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

// Command vier-cli is four in a row in the terminal.
package main

import (
	"os"
	"time"
)

func main() {
	NewSession(os.Stdout, VersusComputer, uint64(time.Now().UnixNano())).Serve(os.Stdin)
}
