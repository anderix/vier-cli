// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Memory is what vier-cli keeps between runs: the last choices, whose turn
// it is to start against the computer, and the game in progress, if any.
type Memory struct {
	Mode        Mode       `json:"mode"`
	Difficulty  Difficulty `json:"difficulty"`
	NextStarter Player     `json:"nextStarter"`
	SavedGame   *SavedGame `json:"savedGame,omitempty"`
}

type SavedGame struct {
	Board Board  `json:"board"`
	Turn  Player `json:"turn"`
}

// statePath is where the memory lives: XDG_STATE_HOME or ~/.local/state on
// Linux, Application Support on macOS, and Local AppData on Windows.
func statePath() (string, error) {
	var base string
	var err error
	switch dir := os.Getenv("XDG_STATE_HOME"); {
	case filepath.IsAbs(dir):
		base = dir
	case runtime.GOOS == "darwin":
		base, err = os.UserConfigDir()
	case runtime.GOOS == "windows":
		base, err = os.UserCacheDir()
	default:
		var home string
		home, err = os.UserHomeDir()
		base = filepath.Join(home, ".local", "state")
	}
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "vier-cli", "state.json"), nil
}

// LoadMemory reads the memory at path. A missing file is a first run, not an error.
func LoadMemory(path string) (Memory, error) {
	memory := Memory{NextStarter: Red}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return memory, nil
	}
	if err != nil {
		return memory, err
	}
	if err := json.Unmarshal(data, &memory); err != nil {
		return Memory{NextStarter: Red}, err
	}
	return memory, nil
}

// Save writes the memory beside path first and renames it into place, so a
// crash part way through leaves the old file whole.
func (m Memory) Save(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	staged := path + ".tmp"
	if err := os.WriteFile(staged, append(data, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(staged, path)
}

func (p Player) MarshalText() ([]byte, error) {
	switch p {
	case Red:
		return []byte("red"), nil
	case Yellow:
		return []byte("yellow"), nil
	default:
		return []byte("nobody"), nil
	}
}

func (p *Player) UnmarshalText(text []byte) error {
	switch string(text) {
	case "red":
		*p = Red
	case "yellow":
		*p = Yellow
	case "nobody":
		*p = Nobody
	default:
		return fmt.Errorf("unknown player %q", text)
	}
	return nil
}

func (m Mode) MarshalText() ([]byte, error) {
	if m == TwoPlayers {
		return []byte("friend"), nil
	}
	return []byte("computer"), nil
}

func (m *Mode) UnmarshalText(text []byte) error {
	switch string(text) {
	case "computer":
		*m = VersusComputer
	case "friend":
		*m = TwoPlayers
	default:
		return fmt.Errorf("unknown mode %q", text)
	}
	return nil
}

func (d Difficulty) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Difficulty) UnmarshalText(text []byte) error {
	level, ok := levels[string(text)]
	if !ok {
		return fmt.Errorf("unknown difficulty %q", text)
	}
	*d = level
	return nil
}

// MarshalJSON writes the board as it is drawn: six rows, top row first.
func (b Board) MarshalJSON() ([]byte, error) {
	rows := make([]string, Rows)
	for row := range Rows {
		for column := range Columns {
			rows[Rows-1-row] += b.cells[row][column].mark()
		}
	}
	return json.Marshal(rows)
}

func (b *Board) UnmarshalJSON(data []byte) error {
	var rows []string
	if err := json.Unmarshal(data, &rows); err != nil {
		return err
	}
	if len(rows) != Rows {
		return fmt.Errorf("board has %d rows, want %d", len(rows), Rows)
	}
	var board Board
	for i, text := range rows {
		if len(text) != Columns {
			return fmt.Errorf("board row %q is not %d cells", text, Columns)
		}
		for column, mark := range text {
			switch mark {
			case 'R':
				board.cells[Rows-1-i][column] = Red
			case 'Y':
				board.cells[Rows-1-i][column] = Yellow
			case '.':
			default:
				return fmt.Errorf("board row %q has %q in it", text, mark)
			}
		}
	}
	*b = board
	return nil
}
