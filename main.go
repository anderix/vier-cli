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

const help = `Type a column number, 1 to 7, to drop a disc there.
  new    start a new game
  show   show the board again
  help   show this list
  quit   leave
`

func run(in io.Reader, out io.Writer) {
	game := NewGame()
	lines := bufio.NewScanner(in)
	fmt.Fprint(out, game.Board)
	for {
		fmt.Fprint(out, prompt(game))
		if !lines.Scan() {
			fmt.Fprintln(out)
			return
		}
		command := strings.ToLower(strings.TrimSpace(lines.Text()))
		switch command {
		case "":
		case "quit":
			return
		case "help":
			fmt.Fprint(out, help)
		case "new":
			game = NewGame()
			fmt.Fprint(out, game.Board)
		case "show":
			fmt.Fprint(out, game.Board)
		default:
			fmt.Fprint(out, play(&game, command))
		}
	}
}

func prompt(game Game) string {
	if game.Over {
		return "> "
	}
	return game.Turn.String() + "> "
}

// play answers a command that is not one of the words in help.
func play(game *Game, command string) string {
	number, err := strconv.Atoi(command)
	if err != nil {
		return fmt.Sprintf("I don't know %q. Type help for the commands.\n", command)
	}
	if err := game.Play(number - 1); err != nil {
		return refusal(err, number) + "\n"
	}
	result := game.Board.String()
	switch {
	case !game.Over:
		return result
	case game.Winner != Nobody:
		return result + game.Winner.String() + " wins! Type new to play again.\n"
	default:
		return result + "The board is full. It's a draw. Type new to play again.\n"
	}
}

func refusal(err error, number int) string {
	switch {
	case errors.Is(err, ErrGameOver):
		return "The game is over. Type new to play again."
	case errors.Is(err, ErrColumnFull):
		return fmt.Sprintf("Column %d is full.", number)
	default:
		return fmt.Sprintf("There is no column %d. The columns are 1 to 7.", number)
	}
}
