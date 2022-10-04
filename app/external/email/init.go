package email

import (
	"bytes"
	"embed"
	"fmt"
	"gamma/app/system/log"
	"html/template"
	"io/fs"
)

type TMPLType string

const (
	html TMPLType = "html"
	text TMPLType = "text"
)

const (
	ResetPassword string = "reset_password"
)

//go:embed templates/*
var templatesFS embed.FS

var templates = map[string]*template.Template{}

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

		log.Infof("%s", template.Name())

		templates[template.Name()] = template
	}
}

func render(key string, renderType TMPLType, data interface{}) (string, error) {
	var tmplOut bytes.Buffer
	// TODO: This probably doesn't work
	if err := templates[fmt.Sprintf("%s_%s", key, renderType)].Execute(&tmplOut, data); err != nil {
		return "", err
	}
	return tmplOut.String(), nil
}
