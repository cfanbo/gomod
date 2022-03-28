package core

import (
	"strings"
)

var (
	Version   = "v0.0.1"
	GitCommit = ""
	BuildMeta = ""
)

// FullVersion returns a version string.
func FullVersion() string {
	var sb strings.Builder
	sb.Grow(len(Version) + len(GitCommit) + len(BuildMeta) + len("-") + len("+"))
	sb.WriteString(Version)
	if BuildMeta != "" {
		sb.WriteString("-" + BuildMeta)
	}
	if GitCommit != "" {
		sb.WriteString("+" + GitCommit)
	}
	return sb.String()
}
