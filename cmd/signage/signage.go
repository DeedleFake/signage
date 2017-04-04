package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/DeedleFake/signage"
)

var (
	tmpl *template.Template
)

func init() {
	tmpl = template.Must(template.New("rss").Parse(`<?xml version='1.0' encoding='UTF-8' ?>
<rss version='2.0'>
	<channel>
		<title>{{ .Type }} Bills</title>

		{{ range .Bills -}}
		<item>
			<title>{{ .Title }}</title>
			<link>{{ .URL }}</link>
			<pubDate>{{ .Date }}</pubDate>
		</item>
		{{- end }}
	</channel>
</rss>`))
}

func handleSigned(rw http.ResponseWriter, req *http.Request) {
	bills, err := signage.GetSigned()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Type":  "Signed",
		"Bills": bills,
	})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(rw, &buf)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handleSigned)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
