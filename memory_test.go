// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func stateFile(t *testing.T) string {
	return filepath.Join(t.TempDir(), "vier-cli", "state.json")
}

// resumed runs a session against the file at path, as a fresh start of the program would.
func resumed(t *testing.T, path, input string) string {
	t.Helper()
	memory, err := LoadMemory(path)
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	NewSession(&out, path, memory, 1).Serve(strings.NewReader(input))
	return out.String()
}

func TestAGameInProgressSurvivesQuitting(t *testing.T) {
	path := stateFile(t)
	resumed(t, path, "friend\n4\n3\nquit\n")
	got := resumed(t, path, "quit\n")
	board := strings.Replace(emptyBoard, " . . . . . . .\n 1", " . . Y R . . .\n 1", 1)
	if want := "Picking up where you left off.\n" + board + "Red> "; got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestAFinishedGameIsNotResumed(t *testing.T) {
	path := stateFile(t)
	resumed(t, path, "friend\n"+moves([]int{1, 1, 2, 2, 3, 3, 4})+"quit\n")
	if got := resumed(t, path, "quit\n"); !strings.HasPrefix(got, twoPlayerGame+emptyBoard) {
		t.Errorf("a finished game came back:\n%s", got)
	}
}

func TestChoicesAreRemembered(t *testing.T) {
	path := stateFile(t)
	resumed(t, path, "hard\n"+moves([]int{1, 1, 1, 1})+"quit\n")
	memory, err := LoadMemory(path)
	if err != nil {
		t.Fatal(err)
	}
	if memory.Mode != VersusComputer || memory.Difficulty != Hard || memory.NextStarter != Red {
		t.Errorf("remembered %+v", memory)
	}
}

func TestTheFileReadsLikeTheBoard(t *testing.T) {
	path := stateFile(t)
	resumed(t, path, "friend\n4\nquit\n")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"mode": "friend"`, `"difficulty": "easy"`, `"...R..."`, `"turn": "yellow"`} {
		if !strings.Contains(string(data), want) {
			t.Errorf("missing %s in:\n%s", want, data)
		}
	}
}

func TestADamagedFileIsRefused(t *testing.T) {
	for _, contents := range []string{
		"not json",
		`{"mode": "solitaire"}`,
		`{"savedGame": {"board": [".......", "......."], "turn": "red"}}`,
		`{"savedGame": {"board": ["...", ".......", ".......", ".......", ".......", "......."], "turn": "red"}}`,
		`{"savedGame": {"board": ["...X...", ".......", ".......", ".......", ".......", "......."], "turn": "red"}}`,
	} {
		path := stateFile(t)
		os.MkdirAll(filepath.Dir(path), 0o755)
		os.WriteFile(path, []byte(contents), 0o644)
		if memory, err := LoadMemory(path); err == nil {
			t.Errorf("%s: accepted as %+v", contents, memory)
		}
	}
}

func TestAMissingFileIsAFirstRun(t *testing.T) {
	memory, err := LoadMemory(stateFile(t))
	if err != nil || memory.Mode != VersusComputer || memory.NextStarter != Red || memory.SavedGame != nil {
		t.Errorf("got %+v, %v", memory, err)
	}
}

func TestAFailedSaveIsReportedOnce(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "file")
	os.WriteFile(blocker, nil, 0o644)
	var out strings.Builder
	NewSession(&out, filepath.Join(blocker, "state.json"), Memory{Mode: TwoPlayers}, 1).Serve(strings.NewReader("4\n4\nquit\n"))
	if got := strings.Count(out.String(), "This game can't be saved"); got != 1 {
		t.Errorf("reported %d times:\n%s", got, out.String())
	}
}
