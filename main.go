// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

// Command vier-cli is four in a row in the terminal.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	run(os.Stdin, os.Stdout)
}

func run(in io.Reader, out io.Writer) {
	var board Board
	lines := bufio.NewScanner(in)
	fmt.Fprint(out, board)
	for {
		fmt.Fprint(out, "> ")
		if !lines.Scan() {
			fmt.Fprintln(out)
			return
		}
		switch strings.TrimSpace(lines.Text()) {
		case "quit":
			return
		default:
			fmt.Fprintln(out, "Type quit to leave.")
		}
	}
}
