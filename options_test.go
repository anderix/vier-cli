// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"errors"
	"flag"
	"strings"
	"testing"
)

func TestOptionsAreRead(t *testing.T) {
	var stderr strings.Builder
	options, err := parseOptions([]string{"--level", "hard", "--seed", "42"}, &stderr)
	if err != nil || options.Level == nil || *options.Level != Hard || options.Seed != 42 || options.Friend {
		t.Errorf("got %+v, %v", options, err)
	}
	for _, args := range [][]string{{"-V"}, {"--version"}} {
		if options, err := parseOptions(args, &stderr); err != nil || !options.Version {
			t.Errorf("%v: got %+v, %v", args, options, err)
		}
	}
}

func TestHelpIsNotAFailure(t *testing.T) {
	for _, args := range [][]string{{"-h"}, {"--help"}} {
		var stderr strings.Builder
		if _, err := parseOptions(args, &stderr); !errors.Is(err, flag.ErrHelp) || stderr.String() != usage {
			t.Errorf("%v: got %v and %q", args, err, stderr.String())
		}
	}
}

func TestBadOptionsAreRefusedWithTheUsage(t *testing.T) {
	for _, c := range []struct {
		args []string
		says string
	}{
		{[]string{"--level", "impossible"}, "choose easy, medium or hard"},
		{[]string{"--seed", "-3"}, "the seed is a whole number"},
		{[]string{"--friend", "--level", "easy"}, "choose one"},
		{[]string{"game.json"}, `only options: "game.json"`},
		{[]string{"--colour"}, "flag provided but not defined: -colour"},
	} {
		var stderr strings.Builder
		_, err := parseOptions(c.args, &stderr)
		if err == nil || !strings.Contains(stderr.String(), c.says) || !strings.HasSuffix(stderr.String(), usage) {
			t.Errorf("%v: got %v and %q", c.args, err, stderr.String())
		}
	}
}

func TestAskingForAGameSetsTheSavedOneAside(t *testing.T) {
	saved := Memory{Mode: TwoPlayers, SavedGame: &SavedGame{Turn: Red}}
	hard := Hard
	if got := (Options{Level: &hard}).apply(saved); got.Mode != VersusComputer || got.Difficulty != Hard || got.SavedGame != nil {
		t.Errorf("--level hard gave %+v", got)
	}
	if got := (Options{Friend: true}).apply(Memory{SavedGame: &SavedGame{}}); got.Mode != TwoPlayers || got.SavedGame != nil {
		t.Errorf("--friend gave %+v", got)
	}
	if got := (Options{}).apply(saved); got.SavedGame == nil {
		t.Error("no options dropped the saved game")
	}
}
