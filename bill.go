package signage

import (
	"bytes"
	"io"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Bill contains information about a specific entry on the White House
// site.
type Bill struct {
	Date  time.Time
	Title string
	URL   string
}

// Summary fetches and scrapes up to length paragraphs of the bill's
// contents. In this instance, 'paragraph' means a single section of
// text, as dilineated by the HTML of the page being scraped. If
// length is less than zero, the whole bill is returned.
func (b Bill) Summary(length int) (string, error) {
	root, err := getHTML(b.URL)
	if err != nil {
		return "", err
	}

	first := findNode(root, func(n *html.Node) bool {
		return (n.Type == html.TextNode) && (strings.HasPrefix(n.Data, "One Hundred"))
	})
	if first == nil {
		return "", nil
	}

	var buf bytes.Buffer
	for i, cur := 0, first.Parent.Parent; ((length < 0) || (i < length)) && (cur != nil); i, cur = i+1, cur.NextSibling {
		writeNode(&buf, cur)
	}

	return buf.String(), nil
}

// writeNode recursively rebuilds HTML source from a node and all of
// its children.
func writeNode(w io.Writer, n *html.Node) {
	var inner func(io.Writer, *html.Node, bool)
	inner = func(w io.Writer, n *html.Node, top bool) {
		if n == nil {
			return
		}

		switch n.Type {
		case html.TextNode:
			io.WriteString(w, n.Data)

		case html.ElementNode:
			io.WriteString(w, "<")
			io.WriteString(w, n.Data)
			if len(n.Attr) > 0 {
				for _, attr := range n.Attr {
					io.WriteString(w, " ")
					io.WriteString(w, attr.Key)
					io.WriteString(w, "='")
					io.WriteString(w, html.EscapeString(attr.Val))
					io.WriteString(w, "'")
				}
			}
			if n.FirstChild == nil {
				io.WriteString(w, " />")
				break
			}
			io.WriteString(w, ">")
			inner(w, n.FirstChild, false)
			io.WriteString(w, "</")
			io.WriteString(w, n.Data)
			io.WriteString(w, ">")
		}

		if !top {
			inner(w, n.NextSibling, false)
		}
	}

	inner(w, n, true)
}
