// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Mode int

const (
	VersusComputer Mode = iota
	TwoPlayers
)

const help = `Type a column number, 1 to 7, to drop a disc there.
  new        start a new game
  computer   start a new game against the computer
  friend     start a new game for two players at this keyboard
  show       show the board again
  help       show this list
  quit       leave
`

// Session is everything that lasts from one game to the next. Against the
// computer you play Red, and the side that starts alternates from game to game.
type Session struct {
	out         io.Writer
	game        Game
	mode        Mode
	nextStarter Player
	rng         SplitMix64
}

func NewSession(out io.Writer, mode Mode, seed uint64) *Session {
	s := &Session{out: out, nextStarter: Red, rng: SplitMix64{state: seed}}
	s.start(mode)
	return s
}

func (s *Session) Serve(in io.Reader) {
	lines := bufio.NewScanner(in)
	for {
		fmt.Fprint(s.out, s.prompt())
		if !lines.Scan() {
			fmt.Fprintln(s.out)
			return
		}
		command := strings.ToLower(strings.TrimSpace(lines.Text()))
		switch command {
		case "":
		case "quit":
			return
		case "help":
			fmt.Fprint(s.out, help)
		case "new":
			s.start(s.mode)
		case "computer":
			s.start(VersusComputer)
		case "friend":
			s.start(TwoPlayers)
		case "show":
			fmt.Fprint(s.out, s.game.Board)
		default:
			s.play(command)
		}
	}
}

func (s *Session) start(mode Mode) {
	s.mode = mode
	starter := Red
	if mode == VersusComputer {
		starter = s.nextStarter
		s.nextStarter = starter.opponent()
	}
	s.game = NewGame(starter)
	fmt.Fprint(s.out, s.game.Board)
	s.computerTurn()
}

func (s *Session) prompt() string {
	if s.game.Over {
		return "> "
	}
	return s.game.Turn.String() + "> "
}

func (s *Session) play(command string) {
	number, err := strconv.Atoi(command)
	if err != nil {
		fmt.Fprintf(s.out, "I don't know %q. Type help for the commands.\n", command)
		return
	}
	if err := s.game.Play(number - 1); err != nil {
		fmt.Fprintln(s.out, refusal(err, number))
		return
	}
	fmt.Fprint(s.out, s.game.Board)
	s.announce()
	s.computerTurn()
}

func (s *Session) computerTurn() {
	if s.mode != VersusComputer || s.game.Over || s.game.Turn != Yellow {
		return
	}
	column := randomMove(s.game.Board, &s.rng)
	s.game.Play(column)
	fmt.Fprintf(s.out, "The computer plays %d.\n", column+1)
	fmt.Fprint(s.out, s.game.Board)
	s.announce()
}

func (s *Session) announce() {
	if !s.game.Over {
		return
	}
	switch {
	case s.game.Winner == Nobody:
		fmt.Fprint(s.out, "The board is full. It's a draw.")
	case s.mode == TwoPlayers:
		fmt.Fprintf(s.out, "%s wins!", s.game.Winner)
	case s.game.Winner == Red:
		fmt.Fprint(s.out, "You win!")
	default:
		fmt.Fprint(s.out, "The computer wins.")
	}
	fmt.Fprintln(s.out, " Type new to play again.")
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
