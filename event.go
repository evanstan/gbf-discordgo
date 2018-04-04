package gbfbot

import (
	"time"
	"github.com/evanstan/go-gbf"
)

const (
	EventCacheTTL = 30 * time.Minute
)

// CachedEvent wraps an event, its details, and its expiration date.
type CachedEvent struct {
	Event     *gbf.Event
	Details   *gbf.EventDetails
	ExpiresAt time.Time
	touchedAt time.Time
}

func (event *CachedEvent) touch() {
	event.touchedAt = time.Now()
}

// CurrentEvents returns the ongoing events.
func (g *GBFBot) CurrentEvents() (events []*CachedEvent, err error) {
	g.eventsMutex.Lock()
	defer g.eventsMutex.Unlock()

	if len(g.currentEvents) > 0 {
		alive := true

		for _, event := range g.currentEvents {
			if time.Now().After(event.ExpiresAt) {
				alive = false
				break
			}
		}

		if alive {
			events = g.currentEvents

			return
		}
	}

	var (
		currentEvents []*gbf.Event
	)

	currentEvents, err = gbf.CurrentEvents()
	if err != nil {
		return
	}

	events = make([]*CachedEvent, len(currentEvents))
	for i, event := range currentEvents {
		cachedEvent := &CachedEvent{
			Event:     event,
			ExpiresAt: time.Now().Add(EventCacheTTL),
		}

		cachedEvent.Details, err = gbf.EventDetailsURL(event.URL)
		if err != nil {
			return
		}

		events[i] = cachedEvent
	}

	g.currentEvents = events

	return
}

// UpcomingEvents returns the upcoming events.
func (g *GBFBot) UpcomingEvents() (events []*CachedEvent, err error) {
	g.eventsMutex.Lock()
	defer g.eventsMutex.Unlock()

	if len(g.upcomingEvents) > 0 {
		alive := true

		for _, event := range g.upcomingEvents {
			if time.Now().After(event.ExpiresAt) {
				alive = false
				break
			}
		}

		if alive {
			events = g.upcomingEvents

			return
		}
	}

	var (
		upcomingEvents []*gbf.Event
	)

	upcomingEvents, err = gbf.UpcomingEvents()
	if err != nil {
		return
	}

	events = make([]*CachedEvent, len(upcomingEvents))
	for i, event := range upcomingEvents {
		cachedEvent := &CachedEvent{
			Event:     event,
			ExpiresAt: time.Now().Add(EventCacheTTL),
		}

		if event.URL != "" {
			cachedEvent.Details, err = gbf.EventDetailsURL(event.URL)
			if err != nil {
				return
			}
		}

		events[i] = cachedEvent
	}

	g.upcomingEvents = events

	return
}
