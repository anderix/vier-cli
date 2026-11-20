// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import "strings"

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
