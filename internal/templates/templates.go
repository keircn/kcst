package templates

import (
	"html/template"
	"io"

	"github.com/keircn/kcst/internal/models"
)

type Templates struct {
	page    *template.Template
	preview *template.Template
}

func New() *Templates {
	return &Templates{
		page:    template.Must(template.New("page").Parse(pageTemplate)),
		preview: template.Must(template.New("preview").Parse(previewTemplate)),
	}
}

func (t *Templates) RenderPage(w io.Writer, data models.PageData) error {
	return t.page.Execute(w, data)
}

func (t *Templates) RenderPreview(w io.Writer, data models.FilePreviewData) error {
	return t.preview.Execute(w, data)
}

const pageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body {
            font-family: monospace;
            max-width: 700px;
            margin: 2rem auto;
            padding: 0 1rem;
            background: #1a1a1a;
            color: #e0e0e0;
            line-height: 1.6;
        }
        h1 {
            color: #fff;
            border-bottom: 1px solid #444;
            padding-bottom: 0.5rem;
        }
        h2 {
            color: #aaa;
            margin-top: 2rem;
        }
        pre {
            background: #2a2a2a;
            padding: 1rem;
            overflow-x: auto;
            border-radius: 4px;
        }
        code {
            background: #2a2a2a;
            padding: 0.2rem 0.4rem;
            border-radius: 2px;
        }
        pre code {
            padding: 0;
        }
        .ascii-art {
            color: #888;
            font-size: 0.85rem;
            line-height: 1.2;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            margin: 1rem 0;
        }
        th, td {
            text-align: left;
            padding: 0.5rem;
            border: 1px solid #444;
        }
        th {
            background: #2a2a2a;
        }
        a {
            color: #6bf;
        }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <p>{{.Message}}</p>

    <h2>Retention Policy</h2>
    <pre class="ascii-art">
min_age  = 3 hours
max_age  = 28 days
max_size = 100 MiB

retention = min_age + (max_age - min_age) * (1 - sqrt(size/max_size))

   days
     28 |.
        | ..
        |   ...
        |      ....
        |          .....
        |               ......
        |                     .......
        |                            ........
      1 |                                    ...............
        +-------------------------------------------------->
        0                    50                          100
                                                         MiB
    </pre>

    <p>Smaller files are retained longer. A 100 MiB file lives ~3 hours, while tiny files can stay up to 28 days.</p>

    <h2>Uploading Files</h2>
    <p>Send a <code>POST</code> request with <code>multipart/form-data</code> containing a <code>file</code> field.</p>

    <table>
        <tr>
            <th>Field</th>
            <th>Description</th>
        </tr>
        <tr>
            <td><code>file</code></td>
            <td>The file to upload (max 100 MiB)</td>
        </tr>
    </table>

    <h2>cURL Examples</h2>
    <pre><code># Upload a file
curl -F 'file=@yourfile.png' {{.BaseURL}}

# Upload from stdin
echo "hello world" | curl -F 'file=@-;filename=hello.txt' {{.BaseURL}}

# Upload with a custom filename
curl -F 'file=@localfile.bin;filename=custom.bin' {{.BaseURL}}</code></pre>

    <table>
        <tr>
            <th>File Size</th>
            <th>Retention</th>
        </tr>
        <tr>
            <td>100 MiB</td>
            <td>~3 hours</td>
        </tr>
        <tr>
            <td>50 MiB</td>
            <td>~9 days</td>
        </tr>
        <tr>
            <td>25 MiB</td>
            <td>~14 days</td>
        </tr>
        <tr>
            <td>10 MiB</td>
            <td>~19 days</td>
        </tr>
        <tr>
            <td>1 MiB</td>
            <td>~25 days</td>
        </tr>
        <tr>
            <td>&lt;1 KiB</td>
            <td>~28 days</td>
        </tr>
    </table>
</body>
</html>
`

const previewTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>

    <meta property="og:title" content="{{.OriginalName}}">
    <meta property="og:description" content="{{.Description}}">
    <meta property="og:url" content="{{.PreviewURL}}">
    <meta property="og:site_name" content="kcst">
{{if eq .MediaType "image"}}
    <meta property="og:type" content="image">
    <meta property="og:image" content="{{.RawURL}}">
    <meta property="og:image:type" content="{{.ContentType}}">
{{else if eq .MediaType "video"}}
    <meta property="og:type" content="video.other">
    <meta property="og:video" content="{{.RawURL}}">
    <meta property="og:video:type" content="{{.ContentType}}">
{{else}}
    <meta property="og:type" content="website">
{{end}}

{{if eq .MediaType "image"}}
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:image" content="{{.RawURL}}">
{{else if eq .MediaType "video"}}
    <meta name="twitter:card" content="player">
    <meta name="twitter:player" content="{{.RawURL}}">
{{else}}
    <meta name="twitter:card" content="summary">
{{end}}
    <meta name="twitter:title" content="{{.OriginalName}}">
    <meta name="twitter:description" content="{{.Description}}">

    <style>
        body {
            font-family: monospace;
            max-width: 800px;
            margin: 2rem auto;
            padding: 0 1rem;
            background: #1a1a1a;
            color: #e0e0e0;
            line-height: 1.6;
        }
        h1 {
            color: #fff;
            border-bottom: 1px solid #444;
            padding-bottom: 0.5rem;
            word-break: break-all;
        }
        .meta {
            background: #2a2a2a;
            padding: 1rem;
            border-radius: 4px;
            margin: 1rem 0;
        }
        .meta-row {
            display: flex;
            padding: 0.25rem 0;
        }
        .meta-label {
            color: #888;
            width: 120px;
            flex-shrink: 0;
        }
        .meta-value {
            color: #e0e0e0;
            word-break: break-all;
        }
        .preview {
            margin: 1.5rem 0;
            text-align: center;
        }
        .preview img, .preview video {
            max-width: 100%;
            max-height: 500px;
            border-radius: 4px;
        }
        a {
            color: #6bf;
        }
        .actions {
            margin-top: 1.5rem;
        }
        .btn {
            display: inline-block;
            background: #333;
            color: #fff;
            padding: 0.5rem 1rem;
            border-radius: 4px;
            text-decoration: none;
            margin-right: 0.5rem;
        }
        .btn:hover {
            background: #444;
        }
    </style>
</head>
<body>
    <h1>{{.OriginalName}}</h1>

    <div class="meta">
        <div class="meta-row">
            <span class="meta-label">Filename:</span>
            <span class="meta-value">{{.Filename}}</span>
        </div>
        <div class="meta-row">
            <span class="meta-label">Size:</span>
            <span class="meta-value">{{.SizeHuman}} ({{.Size}} bytes)</span>
        </div>
        <div class="meta-row">
            <span class="meta-label">Type:</span>
            <span class="meta-value">{{.ContentType}}</span>
        </div>
        <div class="meta-row">
            <span class="meta-label">Uploaded:</span>
            <span class="meta-value">{{.UploadedAt.Format "2006-01-02 15:04:05 UTC"}}</span>
        </div>
        <div class="meta-row">
            <span class="meta-label">Expires:</span>
            <span class="meta-value">{{.ExpiresAt.Format "2006-01-02 15:04:05 UTC"}}</span>
        </div>
    </div>

{{if eq .MediaType "image"}}
    <div class="preview">
        <img src="{{.RawURL}}" alt="{{.OriginalName}}">
    </div>
{{else if eq .MediaType "video"}}
    <div class="preview">
        <video controls>
            <source src="{{.RawURL}}" type="{{.ContentType}}">
            Your browser does not support video playback.
        </video>
    </div>
{{end}}

    <div class="actions">
        <a href="{{.RawURL}}" class="btn">Download / View Raw</a>
        <a href="{{.BaseURL}}" class="btn">Back to Home</a>
    </div>
</body>
</html>
`
