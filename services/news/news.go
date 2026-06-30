package news

import (
	"encoding/json"
	"fmt"
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
	actuals := fetchForexFactoryActuals()
	merged := mergeActuals(events, actuals)

	if len(merged) == 0 {
		log.Println("[news] refresh returned 0 events, keeping previous data")
		return
	}

	sort.Slice(merged, func(i, j int) bool {
		if merged[i].Date == merged[j].Date {
			return merged[i].Time < merged[j].Time
		}
		return merged[i].Date < merged[j].Date
	})

	s.mu.Lock()
	s.events = merged
	s.lastUp = time.Now()
	s.mu.Unlock()
	log.Printf("[news] refreshed %d events (%d actuals matched)", len(merged), len(actuals))
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

type ffEvent struct {
	Name      string `json:"name"`
	Currency  string `json:"currency"`
	Dateline  int64  `json:"dateline"`
	Actual    string `json:"actual"`
	Forecast  string `json:"forecast"`
	Previous  string `json:"previous"`
	Impact    string `json:"impactName"`
}

type ffDay struct {
	Events []ffEvent `json:"events"`
}

type ffComponent struct {
	Days []ffDay `json:"days"`
}

type ffActual struct {
	Date     string
	Time     string
	Currency string
	Actual   string
	Impact   string
	Name     string
}

func extractJSONObject(s string, start int) string {
	if start >= len(s) || s[start] != '{' {
		return ""
	}
	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(s); i++ {
		c := s[i]
		if escaped {
			escaped = false
			continue
		}
		if c == '\\' && inString {
			escaped = true
			continue
		}
		if c == '"' {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		if c == '{' {
			depth++
		} else if c == '}' {
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return ""
}

func fetchForexFactoryActuals() []ffActual {
	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("GET", "https://www.forexfactory.com/calendar?week=thisweek", nil)
	if err != nil {
		log.Printf("[news] forexfactory request error: %v", err)
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[news] forexfactory fetch error: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	html := string(body)

	marker := "window.calendarComponentStates[1] = "
	idx := strings.Index(html, marker)
	if idx < 0 {
		log.Println("[news] could not find calendarComponentStates in forexfactory response")
		return nil
	}
	jsonStart := idx + len(marker)
	jsonStr := extractJSONObject(html, jsonStart)
	if jsonStr == "" {
		log.Println("[news] could not extract calendar JSON object")
		return nil
	}

	var comp ffComponent
	if err := json.Unmarshal([]byte(jsonStr), &comp); err != nil {
		log.Printf("[news] forexfactory JSON parse error: %v", err)
		return nil
	}

	var results []ffActual
	for _, day := range comp.Days {
		for _, ev := range day.Events {
			if ev.Actual == "" || ev.Actual == "-" {
				continue
			}
			t := time.Unix(ev.Dateline, 0).In(iranLoc)
			results = append(results, ffActual{
				Date:     t.Format("2006-01-02"),
				Time:     t.Format("15:04"),
				Currency: ev.Currency,
				Actual:   ev.Actual,
				Impact:   ev.Impact,
				Name:     ev.Name,
			})
		}
	}

	log.Printf("[news] fetched %d actuals from forexfactory.com", len(results))
	return results
}

func mergeActuals(faireconomy []CalendarEvent, actuals []ffActual) []CalendarEvent {
	actualMap := make(map[string]string)
	for _, a := range actuals {
		key := fmt.Sprintf("%s|%s|%s", a.Date, strings.ToUpper(a.Currency), a.Time)
		actualMap[key] = a.Actual
	}

	matched := 0
	for i := range faireconomy {
		timeNorm := faireconomy[i].Time
		if len(timeNorm) >= 4 {
			key := fmt.Sprintf("%s|%s|%s", faireconomy[i].Date, strings.ToUpper(faireconomy[i].Currency), timeNorm)
			if actual, ok := actualMap[key]; ok {
				faireconomy[i].Actual = actual
				matched++
			}
		}
	}

	return faireconomy
}
