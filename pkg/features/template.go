package features

import (
	"context"
	"html/template"
)

// TemplateBindKey is injected into the template root map by view.Renderer.ViewData.
const TemplateBindKey = "_zatrano_feature_eval"

// TemplateFuncMap registers the feature helper for html/template.
// Usage: {{if feature . "my-flag"}} — the first argument must be the template root (dot).
func TemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"feature": featureFromRoot,
	}
}

func featureFromRoot(root any, key string) bool {
	m, _ := root.(map[string]any)
	if m == nil {
		return false
	}
	v, ok := m[TemplateBindKey]
	if !ok || v == nil {
		return false
	}
	e, ok := v.(*Eval)
	if !ok || e == nil {
		return false
	}
	return e.IsEnabled(context.Background(), key)
}
