package signage

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	signed  = "https://www.whitehouse.gov/briefing-room/signed-legislation"
	vetoed  = "https://www.whitehouse.gov/briefing-room/vetoed-legislation"
	pending = "https://www.whitehouse.gov/briefing-room/pending-legislation"
)

func GetSigned() ([]Entry, error) {
	return scrape(signed)
}

func scrape(url string) ([]Entry, error) {
	const (
		DateFormat = "2006-01-02T15:04:05-07:00"
	)

	rsp, err := http.Get(url)
	if err != nil {
		return nil, Error(err)
	}
	defer rsp.Body.Close()

	root, err := html.Parse(rsp.Body)
	if err != nil {
		return nil, Error(err)
	}

	content := findNode(root, func(n *html.Node) bool {
		return (n.Type == html.ElementNode) && (n.Data == "div") && (getAttr(n.Attr, "class") == "view-content")
	})

	var entries []Entry
	for cur := content.FirstChild; cur != nil; cur = cur.NextSibling {
		if (cur.Type != html.ElementNode) || (cur.Data != "div") {
			continue
		}

		found := findNode(cur, func(n *html.Node) bool {
			return (n.Type == html.ElementNode) && (n.Data == "span") && (getAttr(n.Attr, "datatype") == "xsd:dateTime")
		})
		date, err := time.Parse(DateFormat, getAttr(found.Attr, "content"))
		if err != nil {
			return nil, Error(err)
		}

		found = findNode(cur, func(n *html.Node) bool {
			return (n.Type == html.ElementNode) && (n.Data == "h3") && (getAttr(n.Attr, "class") == "field-content")
		})
		title := strings.TrimSpace(found.FirstChild.FirstChild.Data)
		url := getAttr(found.FirstChild.Attr, "href")

		entries = append(entries, Entry{
			Date:  date,
			Title: title,
			URL:   url,
		})
	}

	return entries, nil
}

func findNode(root *html.Node, match func(*html.Node) bool) *html.Node {
	if (root == nil) || match(root) {
		return root
	}

	if found := findNode(root.FirstChild, match); found != nil {
		return found
	}

	return findNode(root.NextSibling, match)
}

func getAttr(attrs []html.Attribute, key string) (val string) {
	for _, attr := range attrs {
		if attr.Key == key {
			return attr.Val
		}
	}

	return
}

type Entry struct {
	Date  time.Time
	Title string
	URL   string
}
