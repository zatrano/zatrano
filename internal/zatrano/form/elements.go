package form

import (
	"fmt"
	"html/template"
	"strings"
)

func (b *Builder) Open() template.HTML {
	actualMethod := "POST"
	if strings.ToUpper(b.method) == "GET" {
		actualMethod = "GET"
	}
	enctype := ""
	if b.isMultipart {
		enctype = ` enctype="multipart/form-data"`
	}
	formTag := fmt.Sprintf(`<form method="%s" action="%s"%s>`, actualMethod, b.action, enctype)
	csrfField := ""
	if b.csrfToken != "" {
		csrfField = fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`, b.csrfField, b.csrfToken)
	}
	methodField := ""
	if m := strings.ToUpper(b.method); m != "GET" && m != "POST" {
		methodField = fmt.Sprintf(`<input type="hidden" name="_method" value="%s">`, m)
	}
	return template.HTML(formTag + "\n" + csrfField + "\n" + methodField)
}

func (b *Builder) Close() template.HTML { return `</form>` }

func (b *Builder) Label(name, text string, attrs ...map[string]string) template.HTML {
	attributes := mergeAttributes(attrs...)
	attributes["for"] = name
	return template.HTML(fmt.Sprintf(`<label %s>%s</label>`, buildAttributes(attributes), text))
}

func (b *Builder) Input(typ, name string, attrs ...map[string]string) template.HTML {
	attributes := mergeAttributes(attrs...)
	finalClass := "form-control"
	if typ == "range" {
		finalClass = "form-range"
	} else if typ == "checkbox" || typ == "radio" {
		finalClass = "form-check-input"
	}

	if b.hasError(name) {
		finalClass += " is-invalid"
	}
	if userClass, ok := attributes["class"]; ok {
		attributes["class"] = userClass
		if !strings.Contains(userClass, "form-control") && !strings.Contains(userClass, "form-range") && !strings.Contains(userClass, "form-check-input") {
			attributes["class"] += " " + finalClass
		}
	} else {
		attributes["class"] = finalClass
	}
	attributes["name"] = name
	attributes["id"] = nameOrID(attributes, name)
	attributes["type"] = typ
	if _, ok := attributes["value"]; !ok {
		value := b.resolveValue(name)
		if value != nil && typ != "password" && typ != "file" {
			attributes["value"] = fmt.Sprintf("%v", value)
		}
	}
	if typ == "password" {
		delete(attributes, "value")
	}
	return template.HTML(fmt.Sprintf(`<input %s>`, buildAttributes(attributes)))
}

func (b *Builder) Text(name string, attrs ...map[string]string) template.HTML { return b.Input("text", name, attrs...) }
func (b *Builder) Email(name string, attrs ...map[string]string) template.HTML { return b.Input("email", name, attrs...) }
func (b *Builder) Password(name string, attrs ...map[string]string) template.HTML { return b.Input("password", name, attrs...) }
func (b *Builder) Hidden(name string, attrs ...map[string]string) template.HTML { return b.Input("hidden", name, attrs...) }
func (b *Builder) File(name string, attrs ...map[string]string) template.HTML { return b.Input("file", name, attrs...) }
func (b *Builder) Number(name string, attrs ...map[string]string) template.HTML { return b.Input("number", name, attrs...) }
func (b *Builder) Date(name string, attrs ...map[string]string) template.HTML { return b.Input("date", name, attrs...) }
func (b *Builder) Time(name string, attrs ...map[string]string) template.HTML { return b.Input("time", name, attrs...) }
func (b *Builder) DatetimeLocal(name string, attrs ...map[string]string) template.HTML { return b.Input("datetime-local", name, attrs...) }
func (b *Builder) Range(name string, attrs ...map[string]string) template.HTML { return b.Input("range", name, attrs...) }

func (b *Builder) Textarea(name string, attrs ...map[string]string) template.HTML {
	attributes := mergeAttributes(attrs...)
	value := b.resolveValue(name)
	delete(attributes, "value")
	finalClass := "form-control"
	if b.hasError(name) { finalClass += " is-invalid" }
	if userClass, ok := attributes["class"]; ok {
		attributes["class"] = userClass + " " + finalClass
	} else {
		attributes["class"] = finalClass
	}
	attributes["name"] = name
	attributes["id"] = nameOrID(attributes, name)
	var valStr string
	if value != nil {
		valStr = fmt.Sprintf("%v", value)
	}
	escapedValue := template.HTMLEscapeString(valStr)
	return template.HTML(fmt.Sprintf(`<textarea %s>%s</textarea>`, buildAttributes(attributes), escapedValue))
}

func (b *Builder) Select(name string, options interface{}, attrs ...map[string]string) template.HTML {
	attributes := mergeAttributes(attrs...)
	selectedValues := b.resolveValueAsSlice(name)
	finalClass := "form-select"
	if b.hasError(name) { finalClass += " is-invalid" }
	if userClass, ok := attributes["class"]; ok {
		attributes["class"] = userClass + " " + finalClass
	} else {
		attributes["class"] = finalClass
	}
	attributes["name"] = name
	attributes["id"] = nameOrID(attributes, name)
	if _, ok := attributes["multiple"]; ok {
		attributes["name"] += "[]"
	}
	optionsHtml := buildOptions(options, selectedValues)
	return template.HTML(fmt.Sprintf(`<select %s>%s</select>`, buildAttributes(attributes), optionsHtml))
}

func (b *Builder) Checkbox(name, value string, attrs ...map[string]string) template.HTML {
	attributes := mergeAttributes(attrs...)
	selectedValue := b.resolveValue(name)
	attributes["type"] = "checkbox"
	attributes["name"] = name
	attributes["id"] = nameOrID(attributes, fmt.Sprintf("%s_%s", name, value))
	attributes["value"] = value
	if isChecked(selectedValue, value) {
		attributes["checked"] = "checked"
	}
	return b.Input("checkbox", name, attributes)
}

func (b *Builder) Radio(name, value string, attrs ...map[string]string) template.HTML {
	attributes := mergeAttributes(attrs...)
	selectedValue := b.resolveValue(name)
	if fmt.Sprintf("%v", selectedValue) == value {
		attributes["checked"] = "checked"
	}
	return b.Input("radio", name, attributes)
}

func (b *Builder) Submit(text string, attrs ...map[string]string) template.HTML {
	attributes := mergeAttributes(attrs...)
	attributes["type"] = "submit"
	if _, ok := attributes["class"]; !ok { attributes["class"] = "btn btn-primary" }
	return template.HTML(fmt.Sprintf(`<button %s>%s</button>`, buildAttributes(attributes), text))
}

func (b *Builder) Button(text string, attrs ...map[string]string) template.HTML {
	attributes := mergeAttributes(attrs...)
	if _, ok := attributes["type"]; !ok { attributes["type"] = "button" }
	if _, ok := attributes["class"]; !ok { attributes["class"] = "btn btn-secondary" }
	return template.HTML(fmt.Sprintf(`<button %s>%s</button>`, buildAttributes(attributes), text))
}

func (b *Builder) FieldError(name string) template.HTML {
	if msg, ok := b.errors[name]; ok {
		return template.HTML(fmt.Sprintf(`<div class="invalid-feedback d-block">%s</div>`, msg))
	}
	return ""
}