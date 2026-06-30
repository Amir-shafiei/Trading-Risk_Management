package news

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type CalendarEvent struct {
	Date         string `json:"date"`
	Time         string `json:"time"`
	Currency     string `json:"currency"`
	EventName    string `json:"event_name"`
	Actual       string `json:"actual"`
	Forecast     string `json:"forecast"`
	Previous     string `json:"previous"`
	Impact       string `json:"impact"`
	ImpactLevel  int    `json:"impact_level"`
	ForecastDiff string `json:"forecast_diff"`
}

type NewsService interface {
	GetLatest() []CalendarEvent
	Refresh()
}

type NewsServiceImpl struct {
	events []CalendarEvent
	mu     sync.RWMutex
	lastUp time.Time
}

var iranLoc *time.Location

func init() {
	var err error
	iranLoc, err = time.LoadLocation("Asia/Tehran")
	if err != nil {
		iranLoc = time.FixedZone("IRST", 3*3600+30*60)
	}
}

func NewNewsService() NewsService {
	s := &NewsServiceImpl{}
	go s.Refresh()
	return s
}

func (s *NewsServiceImpl) Refresh() {
	events := fetchFaireconomy()

	if len(events) == 0 {
		log.Println("[news] refresh returned 0 events, keeping previous data")
		return
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].Date == events[j].Date {
			return events[i].Time < events[j].Time
		}
		return events[i].Date < events[j].Date
	})

	s.mu.Lock()
	s.events = events
	s.lastUp = time.Now()
	s.mu.Unlock()
	log.Printf("[news] refreshed %d events", len(events))
}

func (s *NewsServiceImpl) GetLatest() []CalendarEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if time.Since(s.lastUp) > 30*time.Minute {
		go s.Refresh()
	}

	return s.events
}

func fetchFaireconomy() []CalendarEvent {
	urls := []string{
		"https://nfs.faireconomy.media/ff_calendar_thisweek.json",
		"https://nfs.faireconomy.media/ff_calendar_nextweek.json",
	}

	var allEvents []CalendarEvent
	client := &http.Client{Timeout: 15 * time.Second}

	for _, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("[news] faireconomy fetch error (%s): %v", url, err)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		var raw []struct {
			Title    string `json:"title"`
			Country  string `json:"country"`
			Date     string `json:"date"`
			Impact   string `json:"impact"`
			Forecast string `json:"forecast"`
			Previous string `json:"previous"`
			Actual   string `json:"actual"`
		}
		if err := json.Unmarshal(body, &raw); err != nil {
			continue
		}

		for _, r := range raw {
			t, err := time.Parse(time.RFC3339, r.Date)
			if err != nil {
				continue
			}

			iranTime := t.In(iranLoc)
			impactLevel := 0
			switch strings.ToLower(r.Impact) {
			case "high":
				impactLevel = 3
			case "medium":
				impactLevel = 2
			case "low":
				impactLevel = 1
			}

			allEvents = append(allEvents, CalendarEvent{
				Date:        iranTime.Format("2006-01-02"),
				Time:        iranTime.Format("15:04"),
				Currency:    r.Country,
				EventName:   r.Title,
				Actual:      r.Actual,
				Forecast:    r.Forecast,
				Previous:    r.Previous,
				Impact:      r.Impact,
				ImpactLevel: impactLevel,
			})
		}
	}
	return allEvents
}
