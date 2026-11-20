// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"errors"
	"strings"
)

const (
	Columns = 7
	Rows    = 6
)

type Player int

const (
	Nobody Player = iota
	Red
	Yellow
)

func (p Player) mark() string {
	switch p {
	case Red:
		return "R"
	case Yellow:
		return "Y"
	default:
		return "."
	}
}

// Board holds the discs with row 0 at the bottom, where a dropped disc lands first.
type Board struct {
	cells [Rows][Columns]Player
}

var (
	ErrNoSuchColumn = errors.New("no such column")
	ErrColumnFull   = errors.New("column is full")
)

func (p Player) opponent() Player {
	if p == Red {
		return Yellow
	}
	return Red
}

func (p Player) String() string {
	switch p {
	case Red:
		return "Red"
	case Yellow:
		return "Yellow"
	default:
		return "Nobody"
	}
}

// Drop puts p's disc on top of column, counted from 0, and returns the row it landed in.
func (b *Board) Drop(column int, p Player) (int, error) {
	if column < 0 || column >= Columns {
		return 0, ErrNoSuchColumn
	}
	for row := range Rows {
		if b.cells[row][column] == Nobody {
			b.cells[row][column] = p
			return row, nil
		}
	}
	return 0, ErrColumnFull
}

func (b Board) String() string {
	var out strings.Builder
	for row := Rows - 1; row >= 0; row-- {
		for column := range Columns {
			out.WriteString(" " + b.cells[row][column].mark())
		}
		out.WriteString("\n")
	}
	out.WriteString(" 1 2 3 4 5 6 7\n")
	return out.String()
}
