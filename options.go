// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"time"
)

const usage = `Usage: vier-cli [options]

Four in a row in the terminal, against the computer or a friend.
Once it starts, type help to see the commands.

  --level LEVEL   start a new game against the easy, medium or hard computer
  --friend        start a new game for two players at one keyboard
  --seed N        seed the computer's random choices, to play a game again
  -V, --version   show the version
  -h, --help      show this help
`

// Options is what the command line asked for. Level is nil unless --level was given.
type Options struct {
	Level   *Difficulty
	Friend  bool
	Seed    uint64
	Version bool
}

func parseOptions(args []string, stderr io.Writer) (Options, error) {
	options := Options{Seed: uint64(time.Now().UnixNano())}
	flags := flag.NewFlagSet("vier-cli", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() { fmt.Fprint(stderr, usage) }
	flags.Func("level", "", func(text string) error {
		level, ok := levels[text]
		if !ok {
			return errors.New("choose easy, medium or hard")
		}
		options.Level = &level
		return nil
	})
	flags.BoolVar(&options.Friend, "friend", false, "")
	flags.Func("seed", "", func(text string) error {
		seed, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			return errors.New("the seed is a whole number, 0 or more")
		}
		options.Seed = seed
		return nil
	})
	flags.BoolVar(&options.Version, "version", false, "")
	flags.BoolVar(&options.Version, "V", false, "")
	if err := flags.Parse(args); err != nil {
		return options, err
	}
	if flags.NArg() > 0 {
		return options, fail(flags, "vier-cli takes no file or word, only options: %q", flags.Arg(0))
	}
	if options.Friend && options.Level != nil {
		return options, fail(flags, "--friend and --level each start a game; choose one")
	}
	return options, nil
}

func fail(flags *flag.FlagSet, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	fmt.Fprintln(flags.Output(), err)
	flags.Usage()
	return err
}

// apply turns a request for a new game into the memory a session starts from.
func (o Options) apply(memory Memory) Memory {
	switch {
	case o.Friend:
		memory.Mode = TwoPlayers
		memory.SavedGame = nil
	case o.Level != nil:
		memory.Mode = VersusComputer
		memory.Difficulty = *o.Level
		memory.SavedGame = nil
	}
	return memory
}
