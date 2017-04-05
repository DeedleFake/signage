package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"unicode"

	"github.com/DeedleFake/signage"
)

type marshalFunc func(string, []signage.Bill) (io.Reader, error)

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

type getFunc func() ([]signage.Bill, error)

func getBills(rw http.ResponseWriter, req *http.Request, mode string, get getFunc, marshal marshalFunc) {
	bills, err := get()
	if err != nil {
		log.Printf("Failed to get bills: %v", err)
		http.Error(rw, "Error: Failed to get bills.", http.StatusInternalServerError)
		return
	}

	buf, err := marshal(mode, bills)
	if err != nil {
		log.Printf("Failed to marshal bills: %v", err)
		http.Error(rw, "Error: Failed to marshal bills.", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(rw, buf)
	if err != nil {
		log.Printf("Failed to write to client: %v", err)
		http.Error(rw, "Error: Failed to write to cl... Wait, how can you see this?", http.StatusInternalServerError)
		return
	}
}

var (
	marshallers = map[string]marshalFunc{
		"": marshalRSS,

		".rss":  marshalRSS,
		".json": marshalJSON,
	}

	modes = map[string]getFunc{
		"signed":  signage.GetSigned,
		"vetoed":  signage.GetVetoed,
		"pending": signage.GetPending,
	}
)

func handleList(rw http.ResponseWriter, req *http.Request, root string) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "list", map[string]interface{}{
		"Marshallers": marshallers,
		"Modes":       modes,
		"Root":        root,
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

func mux(root string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("%q request for %q from %v", req.Method, req.URL, req.RemoteAddr)

		name := path.Base(req.URL.Path)
		ext := path.Ext(name)

		mode := name[:len(name)-len(ext)]
		get, ok := modes[mode]
		if !ok {
			handleList(rw, req, root)
			return
		}

		marshal, ok := marshallers[ext]
		if !ok {
			http.Error(rw, fmt.Sprintf("Unknown format: %q", ext), http.StatusBadRequest)
			return
		}

		mode = string(unicode.ToUpper(rune(mode[0]))) + mode[1:]
		getBills(rw, req, mode, get, marshal)
	})
}

func main() {
	addr := flag.String("addr", ":8080", "The address to bind to.")
	root := flag.String("root", "/", "The base URL to build absolute URLs from.")
	flag.Parse()

	log.Fatalln(http.ListenAndServe(*addr, mux(*root)))
}
