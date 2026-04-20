package features

import "testing"

func TestRolloutAllowsStable(t *testing.T) {
	a := rolloutAllows("new-dashboard", 42, 50)
	b := rolloutAllows("new-dashboard", 42, 50)
	if a != b {
		t.Fatalf("same user+key should be stable: %v vs %v", a, b)
	}
}
