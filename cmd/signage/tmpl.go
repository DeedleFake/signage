package main

import (
	"html/template"
	"path"
	"strings"
	"time"
)

var (
	tmpl *template.Template
)

func init() {
	tmpl = template.New("").Funcs(template.FuncMap{
		"rfc2822": func(ts time.Time) string {
			return ts.Format(time.RubyDate)
		},

		"safeHTML": func(str string) template.HTML {
			return template.HTML(str)
		},

		"joinpath": path.Join,
		"join":     strings.Join,
	})

	// TODO: Parallelize summary fetches.
	tmpl = template.Must(tmpl.New("rss").Parse(`{{ "<?xml version='1.0' encoding='UTF-8' ?>" | safeHTML }}
<rss version='2.0'>
	<channel>
		<title>{{ .Type }} Bills</title>

		{{- range .Bills }}
		<item>
			<title>{{ .Title }}</title>
			<link>{{ .URL }}</link>
			<pubDate>{{ .Date | rfc2822 }}</pubDate>
			<description>{{ .Summary -1 }}</description>
		</item>
		{{- end }}
	</channel>
</rss>`))

	tmpl = template.Must(tmpl.New("list").Parse(`<html>
	<head>
		<title>Presidential Bill Lists</title>
	</head>
	<body>
		<ul>
			{{- range $mode, $_ := $.Modes }}
			<li>{{ $mode }}{{ range $format, $_ := $.Marshallers }}{{ if ne $format "" }} (<a href='{{ joinpath $.Root $mode }}{{ $format }}'>{{ $format }}</a>){{ end }}{{ end }}</li>
			{{- end }}
		</ul>
	</body>
</html>`))
}
