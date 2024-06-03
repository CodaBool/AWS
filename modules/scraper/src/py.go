package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

func scrapePY(skipDesc bool) {
	// log := logger.With().Str("func", "scrapePY").Logger()

	defer wg.Done()
	var data []TrendingPY
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "https://hugovk.github.io/top-pypi-packages/top-pypi-packages-30-days.min.json", nil)
	check(err)
	res, err := client.Do(req)
	check(err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	check(err)

	gjson.GetBytes(body, "rows").ForEach(func(index gjson.Result, item gjson.Result) bool {
		if index.Int() == 100 {
			return false
		}
		var project TrendingPY
		item.ForEach(func(key, val gjson.Result) bool {
			if key.String() == "project" {
				project.Name = val.String()
			} else if key.String() == "download_count" {
				project.Downloads = val.Int()
			}
			return true
		})
		data = append(data, project)
		return true
	})

	if skipDesc {
		db.Exec("DELETE FROM trending_pies")
		slog.Info(fmt.Sprintf("+%d short py", len(data)))
		db.Create(data)
	} else {
		slog.Info(fmt.Sprintf("fetching desc for %d .py packages (est. 2 minutes)", len(data)))
		scrapeSummary(data)
	}
}

func scrapeSummary(packages []TrendingPY) {
	// log := logger.With().Str("func", "scrapeSummary").Logger()
	c := colly.NewCollector()
	var newData []TrendingPY
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Delay:      600 * time.Millisecond,
	})
	c.SetRequestTimeout(900 * time.Second)
	c.OnHTML("section", func(e *colly.HTMLElement) {
		summary := e.DOM.Children().Eq(2).Text()
		slog.Debug(fmt.Sprintf("PY, summary = [%s]", strings.Split(summary, "Summary:")))
		summary = strings.Split(summary, "Summary:")[1]
		// slog.Debug(fmt.Sprintf("PY, version = [%s]", strings.Split(summary, "Latest version:")))
		summary = strings.Split(summary, "Latest version:")[0]
		slog.Debug(fmt.Sprintf("PY, package = [%s]", strings.Split(fmt.Sprintf("%v", e.Request.URL), "packages/")))
		packageName := strings.Split(fmt.Sprintf("%v", e.Request.URL), "packages/")[1]
		for _, p := range packages {
			if p.Name == packageName {
				log.Print("+", p.Name)
				newData = append(newData, TrendingPY{
					Description: strings.TrimSpace(summary),
					Name:        p.Name,
					Downloads:   p.Downloads,
				})
			}
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		check(err)
	})
	for _, p := range packages {
		c.Visit("https://www.pypistats.org/packages/" + p.Name)
	}
	slog.Info(fmt.Sprintf("+%d long py", len(newData)))
	db.Exec("DELETE FROM trending_pies")
	db.Create(newData)
}
