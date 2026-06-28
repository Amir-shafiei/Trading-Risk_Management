package news

import (
	"encoding/json"
	"io"
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

func NewNewsService() NewsService {
	s := &NewsServiceImpl{}
	go s.Refresh()
	return s
}

func (s *NewsServiceImpl) Refresh() {
	events := fetchFromAPI()

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
}

func (s *NewsServiceImpl) GetLatest() []CalendarEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if time.Since(s.lastUp) > 30*time.Minute {
		go s.Refresh()
	}

	return s.events
}

func fetchFromAPI() []CalendarEvent {
	urls := []string{
		"https://nfs.faireconomy.media/ff_calendar_thisweek.json",
		"https://nfs.faireconomy.media/ff_calendar_nextweek.json",
	}

	var allEvents []CalendarEvent
	client := &http.Client{Timeout: 15 * time.Second}

	for _, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
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
				Date:        t.Format("2006-01-02"),
				Time:        t.Format("15:04"),
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
