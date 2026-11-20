// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

// searchOrder tries the middle columns first, where more lines of four pass.
var searchOrder = [Columns]int{3, 2, 4, 1, 5, 0, 6}

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
