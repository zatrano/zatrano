package audit

import (
	"encoding/json"
	"testing"
)

func TestDiffJSONPatch_replace(t *testing.T) {
	old := []byte(`{"name":"a","n":1}`)
	new := []byte(`{"name":"b","n":1}`)
	p, err := DiffJSONPatch(old, new)
	if err != nil {
		t.Fatal(err)
	}
	var ops []map[string]any
	if err := json.Unmarshal(p, &ops); err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 || ops[0]["op"] != "replace" {
		t.Fatalf("ops=%v", ops)
	}
}

func TestDiffJSONPatch_add_remove(t *testing.T) {
	old := []byte(`{"a":1}`)
	new := []byte(`{"a":1,"b":2}`)
	p, err := DiffJSONPatch(old, new)
	if err != nil {
		t.Fatal(err)
	}
	var ops []map[string]any
	if err := json.Unmarshal(p, &ops); err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 || ops[0]["op"] != "add" {
		t.Fatalf("ops=%v", ops)
	}
}
