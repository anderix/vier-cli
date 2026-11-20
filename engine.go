// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import "math"

// searchOrder tries the middle columns first, where more lines of four pass,
// so the search meets good moves early and alpha-beta prunes more.
var searchOrder = [Columns]int{3, 2, 4, 1, 5, 0, 6}

const (
	winScore    = 1000
	searchDepth = 4
)

func openColumns(board Board) []int {
	var open []int
	for _, column := range searchOrder {
		if board.CanDrop(column) {
			open = append(open, column)
		}
	}
	return open
}

func randomMove(board Board, rng *SplitMix64) int {
	open := openColumns(board)
	return open[rng.Below(len(open))]
}

// BestMove chooses player's column: a win if there is one, then a block,
// then the best column a search of searchDepth moves can find.
func BestMove(board Board, player Player) int {
	moves := openColumns(board)
	if column, ok := immediateWin(board, player); ok {
		return column
	}
	if column, ok := immediateWin(board, player.opponent()); ok {
		return column
	}
	best := moves[0]
	alpha := -math.MaxInt32
	for _, column := range moves {
		if value := score(board, column, player, searchDepth, 0, -math.MaxInt32, -alpha); value > alpha {
			alpha = value
			best = column
		}
	}
	return best
}

func immediateWin(board Board, player Player) (int, bool) {
	for _, column := range searchOrder {
		next := board
		row, err := next.Drop(column, player)
		if err == nil && next.CompletesFour(Position{column, row}) {
			return column, true
		}
	}
	return 0, false
}

// score is the value to player of dropping into column, seen from player's side.
// A win found sooner scores higher, so the search takes the quickest one.
func score(board Board, column int, player Player, depth, ply, alpha, beta int) int {
	row, _ := board.Drop(column, player)
	if board.CompletesFour(Position{column, row}) {
		return winScore - (ply + 1)
	}
	return -negamax(board, player.opponent(), depth-1, ply+1, -beta, -alpha)
}

func negamax(board Board, player Player, depth, ply, alpha, beta int) int {
	if board.Full() {
		return 0
	}
	if depth == 0 {
		return evaluate(board, player)
	}
	best := -math.MaxInt32
	for _, column := range searchOrder {
		if !board.CanDrop(column) {
			continue
		}
		value := score(board, column, player, depth, ply, alpha, beta)
		best = max(best, value)
		alpha = max(alpha, value)
		if alpha >= beta {
			break
		}
	}
	return best
}

// evaluate guesses how good a position is for player when the search stops:
// discs in the centre column, open threes and twos, and the opponent's open threes.
func evaluate(board Board, player Player) int {
	value := 0
	for row := range Rows {
		if board.cells[row][3] == player {
			value += 3
		}
	}
	for _, line := range lines {
		own, opponent := 0, 0
		for _, cell := range line {
			switch board.at(cell) {
			case player:
				own++
			case Nobody:
			default:
				opponent++
			}
		}
		switch {
		case own == 3 && opponent == 0:
			value += 5
		case own == 2 && opponent == 0:
			value += 2
		case own == 0 && opponent == 3:
			value -= 4
		}
	}
	return value
}
