package gbf

import (
	"net/url"
)

const (
	wikiEventsURLStr = "https://gbf.wiki/Events"
	wikiHomeURLStr   = "https://gbf.wiki/"
)

var (
	WikiHomeURL   *url.URL
	WikiEventsURL *url.URL
)

func init() {
	var (
		err error
	)

	WikiHomeURL, err = url.Parse(wikiHomeURLStr)
	if err != nil {
		panic(err)
	}

	WikiEventsURL, err = url.Parse(wikiEventsURLStr)
	if err != nil {
		panic(err)
	}
}
