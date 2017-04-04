package main

import (
	"html/template"
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
	})

	tmpl = template.Must(tmpl.New("rss").Parse(`{{ "<?xml version='1.0' encoding='UTF-8' ?>" | safeHTML }}
<rss version='2.0'>
	<channel>
		<title>{{ .Type }} Bills</title>

		{{- range .Bills }}
		<item>
			<title>{{ .Title }}</title>
			<link>{{ .URL }}</link>
			<pubDate>{{ .Date | rfc2822 }}</pubDate>
		</item>
		{{- end }}
	</channel>
</rss>`))
}
