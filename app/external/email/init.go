package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
)

// go:embed templates/*
var templatesFS embed.FS
var templates map[string]*template.Template

func init() {
	tmplFiles, err := fs.ReadDir(templatesFS, "templates")
	if err != nil {
		panic(err)
	}

	for _, tmpl := range tmplFiles {
		template, err := template.ParseFS(templatesFS, fmt.Sprintf("templates/%s", tmpl.Name()))
		if err != nil {
			panic(err)
		}

		templates[template.Name()] = template
	}
}

func render(key string, data struct{}) (string, error) {
	var tmplOut bytes.Buffer
	if err := templates[key].Execute(&tmplOut, data); err != nil {
		return "", err
	}
	return tmplOut.String(), nil
}
