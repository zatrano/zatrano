package form

import (
	"fmt"
	"html/template"
	"reflect"
	"sort"
	"strings"
)

type Option struct{ Value, Text string }
type Optgroup struct{ Label string; Options []Option }

func (b *Builder) resolveValue(name string) interface{} {
	cleanName := strings.TrimSuffix(name, "[]")
	if b.oldInput != nil {
		if val, ok := b.oldInput[cleanName]; ok && len(val) > 0 {
			if len(val) == 1 { return val[0] }
			return val
		}
	}
	if b.model != nil {
		return getFieldFromModel(b.model, cleanName)
	}
	return nil
}

func (b *Builder) resolveValueAsSlice(name string) []string {
	value := b.resolveValue(name)
	if value == nil { return nil }
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Slice {
		var result []string
		for i := 0; i < val.Len(); i++ {
			result = append(result, fmt.Sprintf("%v", val.Index(i).Interface()))
		}
		return result
	}
	return []string{fmt.Sprintf("%v", value)}
}

func getFieldFromModel(model interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr { val = val.Elem() }
	if !val.IsValid() || val.Kind() != reflect.Struct { return nil }
	normFieldName := strings.ReplaceAll(strings.Title(strings.ReplaceAll(fieldName, "_", " ")), " ", "")
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("form")
		if tag == "" { tag = field.Tag.Get("json") }
		if strings.Split(tag, ",")[0] == fieldName { return val.Field(i).Interface() }
		if field.Name == normFieldName { return val.Field(i).Interface() }
	}
	return nil
}

func buildAttributes(attrs map[string]string) string {
	var attributes []string
	keys := make([]string, 0, len(attrs))
	for k := range attrs { keys = append(keys, k) }
	sort.Strings(keys)
	for _, k := range keys {
		attributes = append(attributes, fmt.Sprintf(`%s="%s"`, k, template.HTMLEscapeString(attrs[k])))
	}
	return strings.Join(attributes, " ")
}

func mergeAttributes(attrs ...map[string]string) map[string]string {
	merged := make(map[string]string)
	for _, attrMap := range attrs {
		for key, val := range attrMap { merged[key] = val }
	}
	return merged
}

func nameOrID(attrs map[string]string, name string) string {
	if id, ok := attrs["id"]; ok { return id }
	return name
}

func (b *Builder) hasError(name string) bool { _, ok := b.errors[name]; return ok }

func isChecked(selectedValue interface{}, optionValue string) bool {
	if selectedValue == nil { return false }
	val := reflect.ValueOf(selectedValue)
	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			if fmt.Sprintf("%v", val.Index(i).Interface()) == optionValue {
				return true
			}
		}
		return false
	}
	return fmt.Sprintf("%v", selectedValue) == optionValue
}

func buildOptions(options interface{}, selectedValues []string) template.HTML {
	var html strings.Builder
	selectedMap := make(map[string]bool)
	for _, s := range selectedValues { selectedMap[s] = true }
	isSelected := func(val string) string {
		if selectedMap[val] { return " selected" }
		return ""
	}
	switch opts := options.(type) {
	case []Option:
		for _, opt := range opts { html.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`, opt.Value, isSelected(opt.Value), opt.Text)) }
	case []Optgroup:
		for _, group := range opts {
			html.WriteString(fmt.Sprintf(`<optgroup label="%s">`, group.Label))
			for _, opt := range group.Options { html.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`, opt.Value, isSelected(opt.Value), opt.Text)) }
			html.WriteString(`</optgroup>`)
		}
	case map[string]string:
		keys := make([]string, 0, len(opts))
		for k := range opts { keys = append(keys, k) }
		sort.Strings(keys)
		for _, k := range keys { html.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`, k, isSelected(k), opts[k])) }
	}
	return template.HTML(html.String())
}