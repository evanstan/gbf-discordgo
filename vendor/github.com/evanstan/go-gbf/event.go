package gbf

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	eventTimeFormat = "15:04 MST, January _2, 2006"
)

type EventDetails struct {
	ImageURL string
	StartsAt time.Time
	EndsAt   time.Time
}

func EventDetailsURL(eventURLStr string) (details *EventDetails, err error) {
	var (
		eventURL *url.URL
		resp     *http.Response
		doc      *goquery.Document
	)

	eventURL, err = url.Parse(eventURLStr)
	if err != nil {
		err = errors.New("invalid event url")

		return
	}

	resp, err = http.Get(eventURL.String())
	if err != nil {
		return
	}

	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}

	details = &EventDetails{}

	timestamps := doc.Find("#mw-content-text .tooltip .localtime")
	if timestamps.Length() > 1 {
		timestamps.EachWithBreak(func(i int, sel *goquery.Selection) bool {
			if len(sel.Nodes) == 0 || sel.Nodes[0].FirstChild == nil {
				err = errors.New("could not find timestamps node")
				return false
			}

			rawTimestamp := sel.Nodes[0].FirstChild.Data
			rawTimestamp = strings.TrimSpace(rawTimestamp)

			t, err := time.Parse(eventTimeFormat, rawTimestamp)
			if err != nil {
				return false
			}

			if i == 0 {
				details.StartsAt = t
			} else if i == 1 {
				details.EndsAt = t
			}

			return true
		})
	}

	eventImageSel := doc.Find("#mw-content-text > p:first-of-type > img").First()
	if eventImageSel.Length() != 0 {
		var (
			eventImageURL *url.URL
		)

		src := eventImageSel.AttrOr("src", "")
		if src != "" {
			eventImageURL, err = eventURL.Parse(src)
			if err != nil {
				return
			}

			details.ImageURL = eventImageURL.String()
		}
	}

	return
}

type Event struct {
	Title string
	URL   string
}

func CurrentEvents() (events []*Event, err error) {
	var (
		resp *http.Response
		doc  *goquery.Document
	)

	resp, err = http.Get(WikiHomeURL.String())
	if err != nil {
		return
	}

	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}

	td := doc.Find(".wikitable tr:nth-child(8) td:nth-child(1)").First()
	if td.Length() == 0 {
		err = errors.New("events container not found")

		return
	}

	links := td.Find("a")
	events = make([]*Event, 0, links.Length())

	links.Each(func(i int, a *goquery.Selection) {
		var (
			err error

			eventURL *url.URL
		)

		href := a.AttrOr("href", "")
		if href == "" {
			return
		}

		eventURL, err = WikiHomeURL.Parse(href)
		if err != nil {
			return
		}

		title := a.AttrOr("title", "")
		if title == "" {
			return
		}

		event := &Event{
			Title: title,
			URL:   eventURL.String(),
		}

		if !checkDuplicate(events, event) {
			events = append(events, event)
		}
	})

	return
}

func UpcomingEvents() (events []*Event, err error) {
	var (
		resp *http.Response
		doc  *goquery.Document
	)

	resp, err = http.Get(WikiHomeURL.String())
	if err != nil {
		return
	}

	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}

	td := doc.Find(".wikitable tr:nth-child(12) td:nth-child(1)").First()
	if td.Length() == 0 {
		err = errors.New("events container not found")

		return
	}

	linkContainers := td.Find("b")
	events = make([]*Event, 0, linkContainers.Length())

	linkContainers.Each(func(i int, b *goquery.Selection) {
		var (
			err error

			eventURL *url.URL
		)

		title := b.Text()
		title = strings.TrimSpace(title)

		a := b.Find("a").First()
		if a.Length() != 0 {
			href := a.AttrOr("href", "")
			if href == "" {
				return
			}

			eventURL, err = WikiHomeURL.Parse(href)
			if err != nil {
				return
			}
		}

		event := &Event{
			Title: title,
		}

		if eventURL != nil {
			event.URL = eventURL.String()
		}

		if !checkDuplicate(events, event) {
			events = append(events, event)
		}
	})

	return
}

func checkDuplicate(list []*Event, element *Event) bool {
    for _, a := range list {
    	if a.Title == element.Title {
    		return true
    	}
    }
    return false
}
