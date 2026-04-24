package meta

import "runtime/debug"

// Version is a fallback when the binary was not built with module version info
// (e.g. bare `go run` / some CI). Override at link time with:
//
//	-ldflags "-X github.com/zatrano/zatrano/pkg/meta.Version=1.2.3"
var Version = "0.1.0-dev"

// ReportedVersion returns the version users should see.
// After `go install ...@vX` / `@latest` it is the real module version (e.g. v0.0.1);
// for local (devel) builds it falls back to Version.
func ReportedVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return Version
	}
	v := info.Main.Version
	if v == "" || v == "(devel)" {
		return Version
	}
	return v
}
