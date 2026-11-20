// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import "errors"

var ErrGameOver = errors.New("game is over")

// Game is one game from the empty board to its end. Winner stays Nobody in a draw.
type Game struct {
	Board  Board
	Turn   Player
	Winner Player
	Over   bool
}

func NewGame() Game {
	return Game{Turn: Red}
}

// Play drops the current player's disc into column, counted from 0, and passes the turn.
func (g *Game) Play(column int) error {
	if g.Over {
		return ErrGameOver
	}
	if _, err := g.Board.Drop(column, g.Turn); err != nil {
		return err
	}
	g.Winner = g.Board.Winner()
	g.Over = g.Winner != Nobody || g.Board.Full()
	if !g.Over {
		g.Turn = g.Turn.opponent()
	}
	return nil
}
