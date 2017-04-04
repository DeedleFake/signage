package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"

	"github.com/DeedleFake/signage"
)

func marshalRSS(t string, bills []signage.Bill) (io.Reader, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "rss", map[string]interface{}{
		"Type":  t,
		"Bills": bills,
	})
	return &buf, err
}

func marshalJSON(t string, bills []signage.Bill) (io.Reader, error) {
	panic("Not implemented.")
}

func handleSigned(rw http.ResponseWriter, req *http.Request) {
	mode := path.Ext(req.URL.Path)
	if mode == "" {
		mode = ".rss"
	}

	var marshal func(t string, bills []signage.Bill) (io.Reader, error)
	switch mode {
	case ".rss":
		marshal = marshalRSS
	case ".json":
		marshal = marshalJSON
	default:
		http.Error(rw, fmt.Sprintf("Unknown format: %q", mode), http.StatusBadRequest)
		return
	}

	bills, err := signage.GetSigned()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := marshal("Signed", bills)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(rw, buf)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handleSigned)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
