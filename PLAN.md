# Plan

Vier for the terminal: four in a row against a friend or the computer, in Go,
with nothing but the standard library. It reads a line and answers with a line,
so it runs the same in any terminal, over SSH, and on Windows. The rules and the
opponent match vier-desktop and vier-ios: `engine.rs` in vier-desktop is the
specification for the computer player.

Each slice leaves a complete program that builds, passes its tests, and runs.

1. A program that runs: prints an empty board and quits on `quit` or end of input.
2. The board: drop a disc into a column, refuse a full one, show the result.
3. Winning and drawing: two players at one terminal, the game the lux version plays.
4. The command loop: help, new game, bad input answered in plain words.
5. A random opponent, with a seed so a game can be replayed.
6. The search: negamax with alpha-beta, scored the way vier-desktop scores.
7. Difficulty: easy, medium and hard, as vier-desktop defines them.
8. Saving: the game in progress survives quitting, as JSON in the user's state directory.
9. Flags and help text: `--seed`, `--level`, `--version`, `--help`.
10. Release: version stamp, packaging for apt and WinGet, shipped through the fleet's release apps.
