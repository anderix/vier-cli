// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
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
	if got, want := session("quit\nhello\n"), emptyBoard+"> "; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEndOfInputEndsTheSession(t *testing.T) {
	if got, want := session(""), emptyBoard+"> \n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAnythingElseIsAnsweredInWords(t *testing.T) {
	got := session("hello\nquit\n")
	if !strings.Contains(got, "> Type quit to leave.\n> ") {
		t.Errorf("unexpected reply: %q", got)
	}
}
