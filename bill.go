package signage

import (
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

func (b Bill) Summary() ([]Paragraph, error) {
	root, err := getHTML(b.URL)
	if err != nil {
		return nil, err
	}

	first := findNode(root, func(n *html.Node) bool {
		return (n.Type == html.ElementNode) && (n.Data == "p") && (getAttr(n.Attr, "class") == "rtecenter")
	})
	if first == nil {
		return nil, nil
	}

	p := make([]Paragraph, 0, 10)
	for i, cur := 0, first; (i < 10) && (cur != nil); i, cur = i+1, first.NextSibling {
		p = append(p, getParagraph(cur, Paragraph{}))
	}

	return p, nil
}

type Paragraph struct {
	Lines []string
	Type  ParagraphType
}

func getParagraph(n *html.Node, p Paragraph) Paragraph {
	if n == nil {
		return p
	}

	switch n.Type {
	case html.ElementNode:
		switch n.Data {
		case "p":
			if getAttr(n.Attr, "class") == "rtecenter" {
				p.Type |= Centered
			}
			return getParagraph(n.FirstChild, p)

		case "div":
			if strings.Contains(getAttr(n.Attr, "style"), "center;") {
				p.Type |= Centered
			}
			return getParagraph(n.FirstChild, p)

		case "strong":
			p.Type |= Strong
			return getParagraph(n.FirstChild, p)

		case "em":
			p.Type |= Emphasis
			return getParagraph(n.FirstChild, p)
		}

	case html.TextNode:
		if len(strings.TrimSpace(n.Data)) == 0 {
			return getParagraph(n.NextSibling, p)
		}

		for cur := n; cur != nil; cur = cur.NextSibling {
			if cur.Type != html.TextNode {
				continue
			}

			p.Lines = append(p.Lines, cur.Data)
		}
		return p
	}

	return getParagraph(n.FirstChild, p)
}

type ParagraphType uint

const (
	Centered ParagraphType = (1 << iota)
	Strong
	Emphasis
)

func (p ParagraphType) HTML() (before, after string) {
	if p&Emphasis != 0 {
		before = "<em>" + before
		after += "</em>"
	}
	if p&Strong != 0 {
		before = "<strong>" + before
		after += "</string>"
	}

	before = ">" + before
	after += "</p>"
	if p&Centered != 0 {
		before = " style='width:80%;margin:0px auto;'" + before
	}
	before = "<p" + before

	return
}
