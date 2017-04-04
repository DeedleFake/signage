package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"unicode"

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
	buf, err := json.Marshal(map[string]interface{}{
		"type":  t,
		"bills": bills,
	})
	return bytes.NewReader(buf), err
}

func handle(rw http.ResponseWriter, req *http.Request, mode string, get billFunc) {
	format := path.Ext(req.URL.Path)
	if format == "" {
		format = ".rss"
	}

	var marshal func(t string, bills []signage.Bill) (io.Reader, error)
	switch format {
	case ".rss":
		marshal = marshalRSS
	case ".json":
		marshal = marshalJSON
	default:
		http.Error(rw, fmt.Sprintf("Unknown format: %q", format), http.StatusBadRequest)
		return
	}

	bills, err := get()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := marshal(mode, bills)
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

type billFunc func() ([]signage.Bill, error)

var (
	modes = map[string]billFunc{
		"signed":  signage.GetSigned,
		"vetoed":  signage.GetVetoed,
		"pending": signage.GetPending,
	}
)

func mux(rw http.ResponseWriter, req *http.Request) {
	name := path.Base(req.URL.Path)
	ext := path.Ext(name)

	mode := name[:len(name)-len(ext)]
	get, ok := modes[mode]
	if !ok {
		http.Error(rw, fmt.Sprintf("Unknown get: %q", get), http.StatusBadRequest)
		return
	}

	mode = string(unicode.ToUpper(rune(mode[0]))) + mode[1:]
	handle(rw, req, mode, get)
}

func main() {
	log.Fatalln(http.ListenAndServe(":8080", http.HandlerFunc(mux)))
}
