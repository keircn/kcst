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
min_age  = 1 hour
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

    <p>Smaller files are retained longer. A 100 MiB file lives ~1 hour, while tiny files can stay up to 28 days.</p>

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
            <td>~1 hour</td>
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
