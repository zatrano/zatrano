package report_test

import (
	"fmt"
	"testing"

	"github.com/zatrano/framework/core/report"
)

func TestReportCapture(t *testing.T) {
	m := report.New(10)
	m.Capture(fmt.Errorf("boom"), nil)
	m.Capture(fmt.Errorf("again"), nil)
	if m.Count() != 2 {
		t.Fatal(m.Count())
	}
	recent := m.Recent(1)
	if recent[0].Message != "again" {
		t.Fatalf("%+v", recent)
	}
}
