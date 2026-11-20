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
  easy       start a new game against the easy computer
  medium     start a new game against the medium computer
  hard       start a new game against the hard computer
  show       show the board again
  help       show this list
  quit       leave
`

var levels = map[string]Difficulty{"easy": Easy, "medium": Medium, "hard": Hard}

// Session is everything that lasts from one game to the next. Against the
// computer you play Red, and the side that starts alternates from game to game.
// When path is set, the memory is saved there after every change.
type Session struct {
	out    io.Writer
	path   string
	memory Memory
	game   Game
	rng    SplitMix64
}

func NewSession(out io.Writer, path string, memory Memory, seed uint64) *Session {
	s := &Session{out: out, path: path, memory: memory, rng: SplitMix64{state: seed}}
	if saved := memory.SavedGame; saved != nil && saved.Board.Winner() == Nobody && !saved.Board.Full() {
		s.game = NewGame(saved.Turn)
		s.game.Board = saved.Board
		fmt.Fprintln(s.out, "Picking up where you left off.")
		fmt.Fprint(s.out, s.game.Board)
		s.computerTurn()
		s.save()
		return s
	}
	s.start(memory.Mode)
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
			s.start(s.memory.Mode)
		case "computer":
			s.start(VersusComputer)
		case "friend":
			s.start(TwoPlayers)
		case "easy", "medium", "hard":
			s.memory.Difficulty = levels[command]
			s.start(VersusComputer)
		case "show":
			fmt.Fprint(s.out, s.game.Board)
		default:
			s.play(command)
		}
	}
}

func (s *Session) start(mode Mode) {
	s.memory.Mode = mode
	starter := Red
	if mode == VersusComputer {
		if s.memory.NextStarter != Nobody {
			starter = s.memory.NextStarter
		}
		s.memory.NextStarter = starter.opponent()
	}
	s.game = NewGame(starter)
	if mode == VersusComputer {
		fmt.Fprintf(s.out, "A new game against the %s computer.\n", s.memory.Difficulty)
	} else {
		fmt.Fprintln(s.out, "A new game for two players.")
	}
	fmt.Fprint(s.out, s.game.Board)
	s.computerTurn()
	s.save()
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
	s.save()
}

func (s *Session) computerTurn() {
	if s.memory.Mode != VersusComputer || s.game.Over || s.game.Turn != Yellow {
		return
	}
	column := BestMove(s.game.Board, Yellow, s.memory.Difficulty, &s.rng)
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
	case s.memory.Mode == TwoPlayers:
		fmt.Fprintf(s.out, "%s wins!", s.game.Winner)
	case s.game.Winner == Red:
		fmt.Fprint(s.out, "You win!")
	default:
		fmt.Fprint(s.out, "The computer wins.")
	}
	fmt.Fprintln(s.out, " Type new to play again.")
}

// save records the session in its file. A save that fails is reported once
// and not tried again, so a read-only disk does not repeat itself every move.
func (s *Session) save() {
	if s.path == "" {
		return
	}
	s.memory.SavedGame = nil
	if !s.game.Over {
		s.memory.SavedGame = &SavedGame{Board: s.game.Board, Turn: s.game.Turn}
	}
	if err := s.memory.Save(s.path); err != nil {
		fmt.Fprintf(s.out, "This game can't be saved: %v\n", err)
		s.path = ""
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
