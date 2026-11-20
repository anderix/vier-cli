# Security Policy

## Reporting a vulnerability

Please report suspected vulnerabilities privately through GitHub Security Advisories at https://github.com/anderix/vier-cli/security/advisories/new. If you would rather not use GitHub, email david.anderson@excelano.com instead. I aim to respond within seven days.

Please do not open public issues for security problems.

## Supported versions

The latest release receives security fixes. Older versions are not supported. A fix ships in a new release; update the way you installed.

## What vier-cli can access

vier-cli reads what you type and writes to the terminal. It makes no network connections of any kind: no telemetry, no analytics, no update check. It opens no file except its own state file.

## What vier-cli stores

One file, `vier-cli/state.json`, under `$XDG_STATE_HOME` (by default `~/.local/state`) on Linux, Application Support on macOS, and Local AppData on Windows, written with mode `0644` beside a temporary `state.json.tmp` that is renamed over it. It holds the chosen level, whether one or two players are at the board, who moves first next against the computer, and the board and turn of the game in progress. Deleting it forgets all of that and nothing else. vier-cli has no dependencies outside Go's standard library.

## Verifying releases

Every GitHub release includes a `checksums.txt` file listing SHA-256 hashes of all binary archives. Verify any download before running it:

    sha256sum -c checksums.txt --ignore-missing

Release artifacts are built by `build-release`, the release program in excelano/shipping, from a tagged commit using the goreleaser configuration in this repository (`.goreleaser.yml`).
