package audit

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DiffJSONPatch returns an RFC 6902 JSON Patch document (array of operations) between two JSON documents.
// Comparison is shallow on object keys (nested objects are compared by deep JSON string equality of encoded subtrees).
func DiffJSONPatch(oldJSON, newJSON []byte) (json.RawMessage, error) {
	if len(strings.TrimSpace(string(oldJSON))) == 0 {
		oldJSON = []byte("{}")
	}
	if len(strings.TrimSpace(string(newJSON))) == 0 {
		newJSON = []byte("{}")
	}
	var oldObj, newObj map[string]any
	if err := json.Unmarshal(oldJSON, &oldObj); err != nil {
		return nil, fmt.Errorf("audit patch old json: %w", err)
	}
	if err := json.Unmarshal(newJSON, &newObj); err != nil {
		return nil, fmt.Errorf("audit patch new json: %w", err)
	}
	var ops []map[string]any
	for k, nv := range newObj {
		path := "/" + escapeJSONPointerToken(k)
		ov, ok := oldObj[k]
		if !ok {
			ops = append(ops, map[string]any{"op": "add", "path": path, "value": nv})
			continue
		}
		if !jsonEqual(ov, nv) {
			ops = append(ops, map[string]any{"op": "replace", "path": path, "value": nv})
		}
	}
	for k := range oldObj {
		if _, ok := newObj[k]; !ok {
			path := "/" + escapeJSONPointerToken(k)
			ops = append(ops, map[string]any{"op": "remove", "path": path})
		}
	}
	raw, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}

func jsonEqual(a, b any) bool {
	ab, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(ab) == string(bb)
}

func escapeJSONPointerToken(s string) string {
	s = strings.ReplaceAll(s, "~", "~0")
	s = strings.ReplaceAll(s, "/", "~1")
	return s
}
