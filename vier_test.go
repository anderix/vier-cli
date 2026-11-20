// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

const emptyBoard = ` . . . . . . .
 . . . . . . .
 . . . . . . .
 . . . . . . .
 . . . . . . .
 . . . . . . .
 1 2 3 4 5 6 7
`

func TestEmptyBoardShowsSixRowsOfDots(t *testing.T) {
	var board Board
	if got := board.String(); got != emptyBoard {
		t.Errorf("empty board:\n%s\nwant:\n%s", got, emptyBoard)
	}
}

func TestBottomRowPrintsLast(t *testing.T) {
	var board Board
	board.cells[0][0] = Red
	board.cells[Rows-1][Columns-1] = Yellow
	lines := strings.Split(board.String(), "\n")
	if lines[0] != " . . . . . . Y" || lines[Rows-1] != " R . . . . . ." {
		t.Errorf("rows printed in the wrong order:\n%s", board)
	}
}

func session(input string) string {
	var out strings.Builder
	run(strings.NewReader(input), &out)
	return out.String()
}

func TestQuitEndsTheSession(t *testing.T) {
	if got, want := session("quit\nhello\n"), emptyBoard+"Red> "; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEndOfInputEndsTheSession(t *testing.T) {
	if got, want := session(""), emptyBoard+"Red> \n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAnythingElseIsAnsweredInWords(t *testing.T) {
	got := session("hello\nquit\n")
	if !strings.Contains(got, "Red> Type a column from 1 to 7, or quit to leave.\nRed> ") {
		t.Errorf("unexpected reply: %q", got)
	}
}

func TestDiscsStackFromTheBottom(t *testing.T) {
	var board Board
	for want := range Rows {
		row, err := board.Drop(2, Red)
		if err != nil || row != want {
			t.Fatalf("drop %d landed in row %d with %v", want, row, err)
		}
	}
}

func TestAFullColumnRefusesADisc(t *testing.T) {
	var board Board
	for range Rows {
		board.Drop(0, Red)
	}
	before := board
	if _, err := board.Drop(0, Yellow); !errors.Is(err, ErrColumnFull) {
		t.Errorf("got %v, want ErrColumnFull", err)
	}
	if board != before {
		t.Error("a refused drop changed the board")
	}
}

func TestOnlySevenColumnsExist(t *testing.T) {
	var board Board
	for _, column := range []int{-1, Columns} {
		if _, err := board.Drop(column, Red); !errors.Is(err, ErrNoSuchColumn) {
			t.Errorf("column %d: got %v, want ErrNoSuchColumn", column, err)
		}
	}
}

func TestPlayersTakeTurns(t *testing.T) {
	got := session("4\n4\nquit\n")
	if !strings.Contains(got, " . . . Y . . .\n . . . R . . .\n 1 2 3 4 5 6 7\nRed> ") {
		t.Errorf("turns did not alternate:\n%s", got)
	}
}

func TestARefusedMoveKeepsTheTurn(t *testing.T) {
	got := session("9\n1\n1\n1\n1\n1\n1\n1\nquit\n")
	for _, want := range []string{
		"Red> There is no column 9. The columns are 1 to 7.\nRed> ",
		"Red> Column 1 is full.\nRed> ",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

// played drops discs into the given columns, numbered from 1, Red first.
func played(t *testing.T, columns ...int) Board {
	t.Helper()
	var board Board
	turn := Red
	for _, column := range columns {
		if _, err := board.Drop(column-1, turn); err != nil {
			t.Fatalf("column %d: %v", column, err)
		}
		turn = turn.opponent()
	}
	return board
}

func TestThereAre69WinningLines(t *testing.T) {
	if len(lines) != 69 {
		t.Errorf("got %d lines, want 69", len(lines))
	}
}

func TestFourInALineWins(t *testing.T) {
	for _, game := range []struct {
		name    string
		columns []int
		winner  Player
	}{
		{"across", []int{1, 1, 2, 2, 3, 3, 4}, Red},
		{"up", []int{1, 2, 1, 2, 1, 2, 7, 2}, Yellow},
		{"rising diagonal", []int{1, 2, 2, 3, 3, 4, 3, 4, 4, 1, 4}, Red},
		{"falling diagonal", []int{7, 6, 6, 5, 5, 4, 5, 4, 4, 7, 4}, Red},
		{"three is not four", []int{1, 1, 2, 2, 3, 3}, Nobody},
	} {
		if got := played(t, game.columns...).Winner(); got != game.winner {
			t.Errorf("%s: got %s, want %s", game.name, got, game.winner)
		}
	}
}

// drawn fills the board with no four in a line anywhere.
var drawn = []int{
	1, 2, 1, 2, 1, 2, 2, 1, 2, 1, 2, 1,
	3, 4, 3, 4, 3, 4, 4, 3, 4, 3, 4, 3,
	5, 6, 5, 6, 5, 6, 6, 5, 6, 5, 6, 5,
	7, 7, 7, 7, 7, 7,
}

func TestAFullBoardWithNoLineIsADraw(t *testing.T) {
	board := played(t, drawn...)
	if !board.Full() || board.Winner() != Nobody {
		t.Errorf("want a full board with no winner:\n%s", board)
	}
}

func moves(columns []int) string {
	var input strings.Builder
	for _, column := range columns {
		input.WriteString(strconv.Itoa(column) + "\n")
	}
	return input.String()
}

func TestTheGameEndsWithAWinner(t *testing.T) {
	got := session(moves([]int{1, 1, 2, 2, 3, 3, 4, 4}))
	if !strings.HasSuffix(got, " 1 2 3 4 5 6 7\nRed wins!\n") {
		t.Errorf("game did not end on Red's win:\n%s", got)
	}
}

func TestTheGameEndsInADraw(t *testing.T) {
	got := session(moves(drawn))
	if !strings.HasSuffix(got, " 1 2 3 4 5 6 7\nThe board is full. It's a draw.\n") {
		t.Errorf("game did not end in a draw:\n%s", got)
	}
}
