package signage

import (
	"net/http"
	"strings"
	"time"

	"github.com/DeedleFake/signage/errors"

	"golang.org/x/net/html"
)

const (
	base = "https://www.whitehouse.gov"

	signed  = base + "/briefing-room/signed-legislation"
	vetoed  = base + "/briefing-room/vetoed-legislation"
	pending = base + "/briefing-room/pending-legislation"
)

// GetSigned fetches a list of signed bills.
func GetSigned() ([]Bill, error) {
	return scrape(signed)
}

// GetPending fetches a list of pending bills.
func GetPending() ([]Bill, error) {
	return scrape(pending)
}

// GetVetoed fetches a list of vetoed bills.
func GetVetoed() ([]Bill, error) {
	return scrape(vetoed)
}

// scrape pulls a list of bills from a URL.
//
// TODO: Handle scraping multiple pages.
func scrape(url string) ([]Bill, error) {
	const (
		DateFormat = "2006-01-02T15:04:05-07:00"
	)

	rsp, err := http.Get(url)
	if err != nil {
		return nil, errors.Err(err)
	}
	defer rsp.Body.Close()

	root, err := html.Parse(rsp.Body)
	if err != nil {
		return nil, errors.Err(err)
	}

	content := findNode(root, func(n *html.Node) bool {
		return (n.Type == html.ElementNode) && (n.Data == "div") && (getAttr(n.Attr, "class") == "view-content")
	})
	if content == nil {
		return nil, nil
	}

	var entries []Bill
	for cur := content.FirstChild; cur != nil; cur = cur.NextSibling {
		if (cur.Type != html.ElementNode) || (cur.Data != "div") {
			continue
		}

		found := findNode(cur, func(n *html.Node) bool {
			return (n.Type == html.ElementNode) && (n.Data == "span") && (getAttr(n.Attr, "datatype") == "xsd:dateTime")
		})
		date, err := time.Parse(DateFormat, getAttr(found.Attr, "content"))
		if err != nil {
			return nil, errors.Err(err)
		}

		found = findNode(cur, func(n *html.Node) bool {
			return (n.Type == html.ElementNode) && (n.Data == "h3") && (getAttr(n.Attr, "class") == "field-content")
		})
		title := strings.TrimSpace(found.FirstChild.FirstChild.Data)
		url := getAttr(found.FirstChild.Attr, "href")

		entries = append(entries, Bill{
			Date:  date,
			Title: title,
			URL:   base + url,
		})
	}

	return entries, nil
}

// findNode recursively searches an HTML node tree until it finds one
// on which match returns true, at which point it returns that node.
// If no nodes match, it returns nil.
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

// Bill contains information about a specific entry on the White House
// site.
type Bill struct {
	Date  time.Time
	Title string
	URL   string
}
