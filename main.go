// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

// Command vier-cli is four in a row in the terminal.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	run(os.Stdin, os.Stdout)
}

func run(in io.Reader, out io.Writer) {
	var board Board
	turn := Red
	lines := bufio.NewScanner(in)
	fmt.Fprint(out, board)
	for {
		fmt.Fprintf(out, "%s> ", turn)
		if !lines.Scan() {
			fmt.Fprintln(out)
			return
		}
		command := strings.TrimSpace(lines.Text())
		if command == "quit" {
			return
		}
		number, err := strconv.Atoi(command)
		if err != nil {
			fmt.Fprintln(out, "Type a column from 1 to 7, or quit to leave.")
			continue
		}
		if _, err := board.Drop(number-1, turn); err != nil {
			fmt.Fprintln(out, refusal(err, number))
			continue
		}
		fmt.Fprint(out, board)
		turn = turn.opponent()
	}
}

func refusal(err error, number int) string {
	if errors.Is(err, ErrColumnFull) {
		return fmt.Sprintf("Column %d is full.", number)
	}
	return fmt.Sprintf("There is no column %d. The columns are 1 to 7.", number)
}
