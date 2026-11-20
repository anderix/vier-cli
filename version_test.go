// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

import (
	"runtime/debug"
	"testing"
)

func TestPickVersion(t *testing.T) {
	install := &debug.BuildInfo{Main: debug.Module{Version: "v1.4.0"}}
	checkout := &debug.BuildInfo{
		Main:     debug.Module{Version: "1.4.1-0.20260808004201-d7b9d3ff5354+dirty"},
		Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "d7b9d3ff5354"}},
	}
	for _, c := range []struct {
		name    string
		stamped string
		info    *debug.BuildInfo
		want    string
	}{
		{"a release stamp wins", "1.4.0", install, "1.4.0"},
		{"go install reports the module version", versionPlaceholder, install, "1.4.0"},
		{"a checkout keeps the placeholder", versionPlaceholder, checkout, versionPlaceholder},
		{"no build info keeps the placeholder", versionPlaceholder, nil, versionPlaceholder},
		{"devel keeps the placeholder", versionPlaceholder, &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, versionPlaceholder},
	} {
		if got := pickVersion(c.stamped, c.info); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
