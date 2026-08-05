package jsonschema_test

import (
	"testing"

	"github.com/zatrano/framework/core/jsonschema"
)

func TestValidateObject(t *testing.T) {
	schema := jsonschema.Schema{
		Type:     "object",
		Required: []string{"email"},
		Properties: map[string]jsonschema.Schema{
			"email": {Type: "string"},
			"age":   {Type: "integer"},
		},
	}
	errs := jsonschema.Validate(schema, map[string]any{"age": 20})
	if len(errs) == 0 {
		t.Fatal("expected required error")
	}
	if !jsonschema.Valid(schema, map[string]any{"email": "a@b.c", "age": 20}) {
		t.Fatal("expected valid payload")
	}
}

func TestValidateArray(t *testing.T) {
	schema := jsonschema.Schema{
		Type:  "array",
		Items: &jsonschema.Schema{Type: "string"},
	}
	if jsonschema.Valid(schema, []any{"a", 1}) {
		t.Fatal("expected type error on item")
	}
	if !jsonschema.Valid(schema, []any{"a", "b"}) {
		t.Fatal("expected valid array")
	}
}
