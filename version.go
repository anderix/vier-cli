// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Releases set these with -ldflags; see .goreleaser.yml.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const versionPlaceholder = "dev"

func versionLine() string {
	info, _ := debug.ReadBuildInfo()
	v := pickVersion(version, info)
	if commit == "none" || date == "unknown" {
		return "vier-cli " + v
	}
	return fmt.Sprintf("vier-cli %s (%s, %s)", v, commit, date)
}

// pickVersion prefers the release stamp. Without one, a `go install` build
// reports its module version; a build from a checkout carries vcs settings
// and a pseudo-version that reads like a release, so it keeps the placeholder.
func pickVersion(stamped string, info *debug.BuildInfo) string {
	if stamped != versionPlaceholder || info == nil {
		return stamped
	}
	for _, setting := range info.Settings {
		if strings.HasPrefix(setting.Key, "vcs.") {
			return stamped
		}
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return strings.TrimPrefix(v, "v")
	}
	return stamped
}
