package templates

import (
	"html/template"
	"io"

	"github.com/keircn/kcst/internal/models"
)

type Templates struct {
	page *template.Template
}

func New() *Templates {
	return &Templates{
		page: template.Must(template.New("page").Parse(pageTemplate)),
	}
}

func (t *Templates) RenderPage(w io.Writer, data models.PageData) error {
	return t.page.Execute(w, data)
}

const pageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { font-family: system-ui, sans-serif; max-width: 800px; margin: 2rem auto; padding: 0 1rem; }
        h1 { color: #333; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <p>{{.Message}}</p>
</body>
</html>
`
