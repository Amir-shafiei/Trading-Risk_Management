package news

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
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
	forexFacts := fetchFaireconomy()
	forexLiveActuals := scrapeForexLiveActuals()

	events := mergeActuals(forexFacts, forexLiveActuals)

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
	log.Printf("[news] refreshed %d events (%d actuals matched)", len(events), len(forexLiveActuals))
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

type liveActual struct {
	Date     string
	Time     string
	Currency string
	Actual   string
}

func scrapeForexLiveActuals() []liveActual {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium-browser"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://forexfactory.live/"),
		chromedp.WaitVisible("table", chromedp.ByQuery),
		chromedp.Sleep(5*time.Second),
		chromedp.OuterHTML("table", &html, chromedp.ByQuery),
	)
	if err != nil {
		log.Printf("[news] forexfactory.live scrape error: %v", err)
		return nil
	}

	return parseDOMActuals(html)
}

func persianToEnglishDigits(s string) string {
	replacer := strings.NewReplacer(
		"۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4",
		"۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9",
	)
	return replacer.Replace(s)
}

func parseDOMActuals(html string) []liveActual {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil
	}

	var results []liveActual
	currentDate := time.Now().Format("2006-01-02")

	doc.Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")

		if cells.Length() == 1 {
			dayText := strings.TrimSpace(cells.Text())
			if strings.Contains(dayText, "Jun") || strings.Contains(dayText, "Jul") ||
				strings.Contains(dayText, "May") || strings.Contains(dayText, "Aug") ||
				strings.Contains(dayText, "Sep") || strings.Contains(dayText, "Oct") ||
				strings.Contains(dayText, "Nov") || strings.Contains(dayText, "Dec") ||
				strings.Contains(dayText, "Jan") || strings.Contains(dayText, "Feb") ||
				strings.Contains(dayText, "Mar") || strings.Contains(dayText, "Apr") {
				if d := extractEnglishDate(dayText); d != "" {
					currentDate = d
				}
			}
			return
		}

		if cells.Length() < 5 {
			return
		}

		text := strings.TrimSpace(row.Text())
		if strings.Contains(text, "زمان") || strings.Contains(text, "جلسات") {
			return
		}

		timeStr := persianToEnglishDigits(strings.TrimSpace(cells.Eq(0).Text()))
		currency := strings.TrimSpace(cells.Eq(1).Text())
		actualValue := strings.TrimSpace(cells.Eq(4).Text())

		if currency == "" || currency == "-" {
			return
		}

		if actualValue == "-" || actualValue == "" || actualValue == "N/A" {
			return
		}

		timeStr = strings.ReplaceAll(timeStr, " ", "")
		if timeStr == "" || timeStr == "-" {
			return
		}

		results = append(results, liveActual{
			Date:     currentDate,
			Time:     timeStr,
			Currency: currency,
			Actual:   actualValue,
		})
	})

	return results
}

func extractEnglishDate(text string) string {
	months := map[string]string{
		"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04",
		"May": "05", "Jun": "06", "Jul": "07", "Aug": "08",
		"Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12",
	}

	for name, num := range months {
		idx := strings.Index(text, name)
		if idx >= 0 {
			after := text[idx+len(name):]
			after = strings.TrimLeft(after, " /")
			digits := ""
			for _, c := range after {
				if c >= '0' && c <= '9' {
					digits += string(c)
				} else {
					break
				}
			}
			if len(digits) > 0 && len(digits) <= 2 {
				if len(digits) == 1 {
					digits = "0" + digits
				}
				return fmt.Sprintf("%s-%s-%s", time.Now().Year(), num, digits)
			}
		}
	}
	return ""
}

func mergeActuals(faireconomy []CalendarEvent, actuals []liveActual) []CalendarEvent {
	actualMap := make(map[string]string)
	for _, a := range actuals {
		key := fmt.Sprintf("%s|%s|%s", a.Date, strings.ToUpper(a.Currency), a.Time)
		actualMap[key] = a.Actual
	}

	for i := range faireconomy {
		timeNorm := faireconomy[i].Time
		if len(timeNorm) >= 4 {
			key := fmt.Sprintf("%s|%s|%s", faireconomy[i].Date, strings.ToUpper(faireconomy[i].Currency), timeNorm)
			if actual, ok := actualMap[key]; ok {
				faireconomy[i].Actual = actual
				continue
			}
		}
	}

	return faireconomy
}
