// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"errors"
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
