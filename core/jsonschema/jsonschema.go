package jsonschema

import (
	"fmt"
	"reflect"
	"strings"
)

// Schema is a minimal JSON Schema subset for request validation stubs.
type Schema struct {
	Type       string            `json:"type,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Required   []string          `json:"required,omitempty"`
	Items      *Schema           `json:"items,omitempty"`
}

// ValidationError describes a single schema violation.
type ValidationError struct {
	Path    string
	Message string
}

func (e ValidationError) Error() string {
	if e.Path == "" {
		return e.Message
	}
	return e.Path + ": " + e.Message
}

// Validate checks data against schema and returns all errors.
func Validate(schema Schema, data any) []error {
	var errs []error
	validate("", schema, data, &errs)
	return errs
}

// Valid reports whether data matches schema.
func Valid(schema Schema, data any) bool {
	return len(Validate(schema, data)) == 0
}

func validate(path string, schema Schema, data any, errs *[]error) {
	if schema.Type != "" && !typeMatches(schema.Type, data) {
		*errs = append(*errs, ValidationError{Path: path, Message: "expected type " + schema.Type})
		return
	}

	switch strings.ToLower(schema.Type) {
	case "object", "":
		obj, ok := asMap(data)
		if schema.Type == "object" && !ok {
			return
		}
		if ok {
			for _, key := range schema.Required {
				if _, exists := obj[key]; !exists {
					*errs = append(*errs, ValidationError{
						Path:    join(path, key),
						Message: "required",
					})
				}
			}
			for key, prop := range schema.Properties {
				if v, exists := obj[key]; exists {
					validate(join(path, key), prop, v, errs)
				}
			}
		}
	case "array":
		rv := reflect.ValueOf(data)
		if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
			return
		}
		if schema.Items != nil {
			for i := 0; i < rv.Len(); i++ {
				validate(fmt.Sprintf("%s[%d]", path, i), *schema.Items, rv.Index(i).Interface(), errs)
			}
		}
	}
}

func typeMatches(want string, data any) bool {
	if data == nil {
		return strings.EqualFold(want, "null")
	}
	switch strings.ToLower(want) {
	case "object":
		_, ok := asMap(data)
		return ok
	case "array":
		k := reflect.ValueOf(data).Kind()
		return k == reflect.Slice || k == reflect.Array
	case "string":
		_, ok := data.(string)
		return ok
	case "number", "integer":
		switch data.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			if strings.EqualFold(want, "integer") {
				switch v := data.(type) {
				case float64:
					return v == float64(int64(v))
				case float32:
					return v == float32(int64(v))
				}
			}
			return true
		default:
			return false
		}
	case "boolean":
		_, ok := data.(bool)
		return ok
	case "null":
		return data == nil
	default:
		return true
	}
}

func asMap(data any) (map[string]any, bool) {
	switch v := data.(type) {
	case map[string]any:
		return v, true
	default:
		return nil, false
	}
}

func join(path, key string) string {
	if path == "" {
		return key
	}
	return path + "." + key
}
