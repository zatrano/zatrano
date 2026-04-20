package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ViewOptions controls what files the View generator creates.
type ViewOptions struct {
	// Layout is the layout template the generated views extend (default: "layouts/app").
	Layout string
	// WithForm creates an additional create.html / edit.html with form scaffolding.
	WithForm bool
	// DryRun reports what would be created without writing any files.
	DryRun bool
}

// View scaffolds server-rendered HTML view templates for a CRUD module under
// viewsRoot/<name>/.
//
// Files generated:
//
//	views/<name>/index.html   — list page
//	views/<name>/show.html    — detail page
//	views/<name>/create.html  — new-record form  (WithForm)
//	views/<name>/edit.html    — edit-record form (WithForm)
func View(viewsRoot, rawName string, opts ViewOptions) ([]string, error) {
	name := normalizeName(rawName)
	if name == "" {
		return nil, fmt.Errorf("invalid view name %q (use letters, digits, _ or -)", rawName)
	}
	if opts.Layout == "" {
		opts.Layout = "layouts/app"
	}

	pascal := snakeToPascal(name)
	base := filepath.Join(viewsRoot, name)

	files := map[string]string{
		"index.html": tmplViewIndex(name, pascal, opts.Layout),
		"show.html":  tmplViewShow(name, pascal, opts.Layout),
	}
	if opts.WithForm {
		files["create.html"] = tmplViewCreate(name, pascal, opts.Layout)
		files["edit.html"] = tmplViewEdit(name, pascal, opts.Layout)
	}

	var written []string
	for fn, body := range files {
		path := filepath.Join(base, fn)
		written = append(written, path)
		if opts.DryRun {
			continue
		}
		if err := os.MkdirAll(base, 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return nil, err
		}
	}
	return written, nil
}

// ----------------------------------------------------------------------------
// View templates
// ----------------------------------------------------------------------------

func tmplViewIndex(name, pascal, layout string) string {
	title := strings.ReplaceAll(strings.Title(strings.ReplaceAll(name, "_", " ")), "-", " ") //nolint:staticcheck
	return fmt.Sprintf(`{{extends "%s"}}

{{block "title"}}%s List{{end}}

{{block "content"}}
<div class="container">
  <div class="page-header">
    <h1>%s</h1>
    <a href="/%s/create" class="btn btn-primary">New %s</a>
  </div>

  {{if .Items}}
  <table class="table">
    <thead>
      <tr>
        <th>ID</th>
        <th>Created At</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {{range .Items}}
      <tr>
        <td>{{.ID}}</td>
        <td>{{.CreatedAt}}</td>
        <td>
          <a href="/%s/{{.ID}}" class="btn btn-sm btn-secondary">View</a>
          <a href="/%s/{{.ID}}/edit" class="btn btn-sm btn-primary">Edit</a>
        </td>
      </tr>
      {{end}}
    </tbody>
  </table>
  {{else}}
  <p class="empty-state">No %s found.</p>
  {{end}}
</div>
{{end}}
`, layout, title, title, name, pascal, name, name, title)
}

func tmplViewShow(name, pascal, layout string) string {
	title := strings.ReplaceAll(strings.Title(strings.ReplaceAll(name, "_", " ")), "-", " ") //nolint:staticcheck
	return fmt.Sprintf(`{{extends "%s"}}

{{block "title"}}%s Detail{{end}}

{{block "content"}}
<div class="container">
  <div class="page-header">
    <h1>%s #{{.Item.ID}}</h1>
    <a href="/%s" class="btn btn-secondary">Back</a>
    <a href="/%s/{{.Item.ID}}/edit" class="btn btn-primary">Edit</a>
  </div>

  <div class="card">
    <div class="card-body">
      {{/* TODO: render item fields */}}
      <pre>{{json .Item}}</pre>
    </div>
  </div>
</div>
{{end}}
`, layout, title, title, name, name)
}

func tmplViewCreate(name, pascal, layout string) string {
	title := strings.ReplaceAll(strings.Title(strings.ReplaceAll(name, "_", " ")), "-", " ") //nolint:staticcheck
	return fmt.Sprintf(`{{extends "%s"}}

{{block "title"}}New %s{{end}}

{{block "content"}}
<div class="container">
  <div class="page-header">
    <h1>New %s</h1>
    <a href="/%s" class="btn btn-secondary">Cancel</a>
  </div>

  {{form_open "/%s" "POST"}}
    {{csrf_field .CSRF}}

    {{/* Example field — replace with your actual fields */}}
    {{template "components/form-input" (dict
      "Type"  "text"
      "Name"  "name"
      "Label" "Name"
      "Value" (old "name" .Old)
      "Required" true
      "Error" (index .Errors "name")
    )}}

    {{template "components/button" (dict "Label" "Create %s" "Type" "submit" "Variant" "primary")}}
  {{form_close}}
</div>
{{end}}
`, layout, title, title, name, name, pascal)
}

func tmplViewEdit(name, pascal, layout string) string {
	title := strings.ReplaceAll(strings.Title(strings.ReplaceAll(name, "_", " ")), "-", " ") //nolint:staticcheck
	return fmt.Sprintf(`{{extends "%s"}}

{{block "title"}}Edit %s{{end}}

{{block "content"}}
<div class="container">
  <div class="page-header">
    <h1>Edit %s #{{.Item.ID}}</h1>
    <a href="/%s/{{.Item.ID}}" class="btn btn-secondary">Cancel</a>
  </div>

  {{form_open (printf "/%s/%%v" .Item.ID) "POST"}}
    {{csrf_field .CSRF}}
    <input type="hidden" name="_method" value="PUT">

    {{/* Example field — replace with your actual fields */}}
    {{template "components/form-input" (dict
      "Type"  "text"
      "Name"  "name"
      "Label" "Name"
      "Value" (old "name" .Old | default .Item.Name)
      "Required" true
      "Error" (index .Errors "name")
    )}}

    {{template "components/button" (dict "Label" "Update %s" "Type" "submit" "Variant" "primary")}}
  {{form_close}}
</div>
{{end}}
`, layout, title, title, name, name, pascal)
}
