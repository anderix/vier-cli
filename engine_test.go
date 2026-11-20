// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"slices"
	"strings"
	"testing"
	"time"
)

// picture builds a board from six rows of R, Y and dots, top row first.
func picture(t *testing.T, rows string) Board {
	t.Helper()
	lines := strings.Split(rows, "\n")
	if len(lines) != Rows {
		t.Fatalf("picture has %d rows", len(lines))
	}
	var board Board
	for i, line := range lines {
		for column, cell := range line {
			switch cell {
			case 'R':
				board.cells[Rows-1-i][column] = Red
			case 'Y':
				board.cells[Rows-1-i][column] = Yellow
			}
		}
	}
	return board
}

func TestTakesAnImmediateWin(t *testing.T) {
	for _, c := range []struct {
		name    string
		board   string
		player  Player
		columns []int
	}{
		{"horizontal", ".......\n.......\n.......\n.......\n..RR...\n.YYY.R.", Yellow, []int{0, 4}},
		{"vertical", ".......\n.......\n.......\nR......\nR...Y..\nR..YY..", Red, []int{0}},
		{"diagonal", ".......\n.......\n.......\n..YR...\n.YRY...\nYRYR...", Yellow, []int{3}},
	} {
		if got := BestMove(picture(t, c.board), c.player); !slices.Contains(c.columns, got) {
			t.Errorf("%s: played %d, want one of %v", c.name, got, c.columns)
		}
	}
}

func TestBlocksAnImmediateWin(t *testing.T) {
	for _, c := range []struct {
		name   string
		board  string
		column int
	}{
		{"horizontal", ".......\n.......\n.......\n.......\n...Y...\nYRRR.Y.", 4},
		{"vertical", ".......\n.......\n.......\n..R....\n..R.Y..\n..R.YY.", 2},
		{"diagonal", ".......\n.......\n.......\n..RY...\n.RYR...\nRYRY...", 3},
	} {
		if got := BestMove(picture(t, c.board), Yellow); got != c.column {
			t.Errorf("%s: played %d, want %d", c.name, got, c.column)
		}
	}
}

func TestAnOwnWinBeatsABlock(t *testing.T) {
	board := picture(t, ".......\n.......\n.......\n......Y\n......Y\n.RRR..Y")
	if got := BestMove(board, Yellow); got != 6 {
		t.Errorf("played %d, want 6", got)
	}
}

func TestPlaysTheOnlyOpenColumn(t *testing.T) {
	board := picture(t, "RYRY.YR\nRYRY.YR\nYRYRRRY\nYRYRYRY\nRYRYRYR\nRYRYRYR")
	if got := BestMove(board, Red); got != 4 {
		t.Errorf("played %d, want 4", got)
	}
}

func TestTheSearchNeverLosesAGameAgainstRandomPlay(t *testing.T) {
	rng := SplitMix64{state: 11}
	for range 20 {
		game := NewGame(Red)
		for !game.Over {
			column := randomMove(game.Board, &rng)
			if game.Turn == Yellow {
				column = BestMove(game.Board, Yellow)
			}
			if err := game.Play(column); err != nil {
				t.Fatalf("column %d: %v", column+1, err)
			}
		}
		if game.Winner == Red {
			t.Fatalf("random play beat the search:\n%s", game.Board)
		}
	}
}

func TestTheSearchIsFast(t *testing.T) {
	for _, board := range []Board{{}, played(t, 4, 4, 3, 5, 5, 3, 6, 2, 2, 6)} {
		started := time.Now()
		BestMove(board, Red)
		if elapsed := time.Since(started); elapsed > 100*time.Millisecond {
			t.Errorf("took %v", elapsed)
		}
	}
}
