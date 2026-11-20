# vier-cli

Four in a row in the terminal: drop discs into a seven-by-six board and line up four across, down or diagonally before the other side does. You play red against the computer, or hand the keyboard back and forth with a friend. vier-cli reads a line and answers with a line, with no full-screen takeover, so it plays the same in any terminal, over SSH, and in the plain Windows console. It has no account, no network and no dependencies beyond Go's standard library.

It is the terminal sibling of [Vier for the Linux desktop](https://github.com/anderix/vier-desktop) and [Vier for iPhone and Apple Watch](https://github.com/anderix/vier-ios). All of them play by the same rules against the same opponent: from the same position and the same random numbers, vier-cli and the desktop game choose the same move. Nothing passes between them.

vier-cli is also the program built, slice by slice, in [Outgrowing Lux](https://github.com/anderix/outgrowing-lux), a book that starts from a lux version of the game and grows it into this one.

## Install

On Debian and Ubuntu, from the Excelano apt repository, so `apt upgrade` keeps it current:

```sh
curl -fsSL https://excelano.com/apt/setup.sh | sudo sh
sudo apt install vier-cli
```

On Windows, unzip `vier-cli.exe` from the [latest release](https://github.com/anderix/vier-cli/releases/latest) into a folder on your `PATH`.

Anywhere with a Go toolchain:

```sh
go install github.com/anderix/vier-cli@latest
```

## Playing

Run `vier-cli` and type a column number to drop your disc there. The board is drawn after every move, with red as `R`, yellow as `Y` and the column numbers underneath:

```
 . . . . . . .
 . . . . . . .
 . . . . . . .
 . . . . . . .
 . . . Y . . .
 . . R R . . .
 1 2 3 4 5 6 7
Red>
```

Against the computer, who moves first alternates from game to game. Easy, Medium and Hard set how far ahead it looks and how often it plays a random move instead; every level takes a winning move when it has one, and Medium and Hard always block yours. Type `help` at the prompt for the commands that start a new game, change the level or switch to two players, and `vier-cli --help` for the options that do the same from the command line. `--seed` fixes the computer's random choices, so a game can be played again move for move.

Quitting keeps the game in progress, and the next run picks it up. vier-cli keeps that game, the chosen level, and who moves first next in `vier-cli/state.json` under `$XDG_STATE_HOME` (by default `~/.local/state`) on Linux, Application Support on macOS, and Local AppData on Windows.

## Building

`go build` builds the program and `go test ./...` runs the tests. `go.mod` names the Go version it needs.
