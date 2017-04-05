package signage

import (
	"time"
)

// Bill contains information about a specific entry on the White House
// site.
type Bill struct {
	Date  time.Time
	Title string
	URL   string
}
