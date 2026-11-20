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

var difficulties = []Difficulty{Easy, Medium, Hard}

func pick(board Board, player Player, difficulty Difficulty, seed uint64) int {
	return BestMove(board, player, difficulty, &SplitMix64{state: seed})
}

func TestEveryLevelTakesAnImmediateWin(t *testing.T) {
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
		for _, difficulty := range difficulties {
			for seed := range uint64(40) {
				if got := pick(picture(t, c.board), c.player, difficulty, seed); !slices.Contains(c.columns, got) {
					t.Errorf("%s at %s, seed %d: played %d, want one of %v", c.name, difficulty, seed, got, c.columns)
				}
			}
		}
	}
}

func TestMediumAndHardBlockAnImmediateWin(t *testing.T) {
	for _, c := range []struct {
		name   string
		board  string
		column int
	}{
		{"horizontal", ".......\n.......\n.......\n.......\n...Y...\nYRRR.Y.", 4},
		{"vertical", ".......\n.......\n.......\n..R....\n..R.Y..\n..R.YY.", 2},
		{"diagonal", ".......\n.......\n.......\n..RY...\n.RYR...\nRYRY...", 3},
	} {
		for _, difficulty := range []Difficulty{Medium, Hard} {
			if got := pick(picture(t, c.board), Yellow, difficulty, 1); got != c.column {
				t.Errorf("%s at %s: played %d, want %d", c.name, difficulty, got, c.column)
			}
		}
	}
}

func TestAnOwnWinBeatsABlock(t *testing.T) {
	board := picture(t, ".......\n.......\n.......\n......Y\n......Y\n.RRR..Y")
	for _, difficulty := range difficulties {
		if got := pick(board, Yellow, difficulty, 1); got != 6 {
			t.Errorf("%s played %d, want 6", difficulty, got)
		}
	}
}

func TestPlaysTheOnlyOpenColumn(t *testing.T) {
	board := picture(t, "RYRY.YR\nRYRY.YR\nYRYRRRY\nYRYRYRY\nRYRYRYR\nRYRYRYR")
	for _, difficulty := range difficulties {
		if got := pick(board, Red, difficulty, 1); got != 4 {
			t.Errorf("%s played %d, want 4", difficulty, got)
		}
	}
}

func TestHardNeverLosesAGameAgainstRandomPlay(t *testing.T) {
	rng := SplitMix64{state: 11}
	for range 20 {
		game := NewGame(Red)
		for !game.Over {
			open := openColumns(game.Board)
			column := open[rng.Below(len(open))]
			if game.Turn == Yellow {
				column = BestMove(game.Board, Yellow, Hard, &rng)
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
		pick(board, Red, Hard, 1)
		if elapsed := time.Since(started); elapsed > 100*time.Millisecond {
			t.Errorf("took %v", elapsed)
		}
	}
}
