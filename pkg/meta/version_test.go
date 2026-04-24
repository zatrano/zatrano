package meta

import (
	"runtime/debug"
	"testing"
)

func TestReportedVersion_nonEmpty(t *testing.T) {
	v := ReportedVersion()
	if v == "" {
		t.Fatal("ReportedVersion() is empty")
	}
	// Under `go test`, ReadBuildInfo usually reports (devel); expect fallback Version.
	if _, ok := debug.ReadBuildInfo(); ok {
		_ = v
	}
}
