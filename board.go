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

type Position struct {
	Column, Row int
}

type Line [4]Position

// lines is every run of four cells on the board: across, up, and both diagonals.
var lines = winningLines()

// linesThrough lists, for each cell, the lines of four that pass through it.
var linesThrough = func() (table [Rows][Columns][]Line) {
	for _, line := range lines {
		for _, cell := range line {
			table[cell.Row][cell.Column] = append(table[cell.Row][cell.Column], line)
		}
	}
	return table
}()

func winningLines() []Line {
	directions := []Position{{1, 0}, {0, 1}, {1, 1}, {1, -1}}
	var all []Line
	for column := range Columns {
		for row := range Rows {
			for _, step := range directions {
				end := Position{column + 3*step.Column, row + 3*step.Row}
				if end.Column < 0 || end.Column >= Columns || end.Row < 0 || end.Row >= Rows {
					continue
				}
				var line Line
				for i := range line {
					line[i] = Position{column + i*step.Column, row + i*step.Row}
				}
				all = append(all, line)
			}
		}
	}
	return all
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

func (b Board) CanDrop(column int) bool {
	return column >= 0 && column < Columns && b.cells[Rows-1][column] == Nobody
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

func (b Board) at(p Position) Player {
	return b.cells[p.Row][p.Column]
}

// Winner returns the player with four in a line, or Nobody.
func (b Board) Winner() Player {
	for _, line := range lines {
		first := b.at(line[0])
		if first != Nobody && b.at(line[1]) == first && b.at(line[2]) == first && b.at(line[3]) == first {
			return first
		}
	}
	return Nobody
}

// CompletesFour reports whether the disc at p is part of four in a line.
func (b Board) CompletesFour(p Position) bool {
	player := b.at(p)
	if player == Nobody {
		return false
	}
	for _, line := range linesThrough[p.Row][p.Column] {
		if b.at(line[0]) == player && b.at(line[1]) == player && b.at(line[2]) == player && b.at(line[3]) == player {
			return true
		}
	}
	return false
}

func (b Board) Full() bool {
	for column := range Columns {
		if b.CanDrop(column) {
			return false
		}
	}
	return true
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
