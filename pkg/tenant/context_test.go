package tenant

import (
	"context"
	"testing"
)

func TestWithFromContext(t *testing.T) {
	ctx := context.Background()
	ctx = WithContext(ctx, Info{Key: "acme", NumericID: 0})
	info, ok := FromContext(ctx)
	if !ok || info.Key != "acme" {
		t.Fatalf("got ok=%v info=%+v", ok, info)
	}
}

func TestParseNumericKey(t *testing.T) {
	if ParseNumericKey("42") != 42 {
		t.Fatal()
	}
	if ParseNumericKey("x") != 0 {
		t.Fatal()
	}
}
