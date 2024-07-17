package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
var extraText = regexp.MustCompile(`Includes \d+ games`)

func scrapeGames() {
	defer wg.Done()
	var data []TrendingGame
	c := colly.NewCollector()
	c.OnHTML("#search_resultsRows", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, e *colly.HTMLElement) {
			title := e.DOM.Find(".search_name").Text()
			match := extraText.FindString(title)
			if match != "" {
				title = strings.ReplaceAll(title, match, "")
			}
			title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
			title = strings.TrimSpace(strings.ReplaceAll(title, "VR Supported", ""))
			title = nonAlphanumericRegex.ReplaceAllString(title, "")

			msrp := e.DOM.Find(".discount_original_price").Text()
			msrp = strings.TrimSpace(msrp)

			price := e.DOM.Find(".discount_final_price").Text()
			price = strings.TrimSpace(price)

			// log.Println("price", price, "title", title, "msrp", msrp)

			if msrp == "" {
				msrp = price
				if price == "" {
					// log.Println("")
					msrp = "Free"
					price = "Free"
				}
			}

			if price != "Free" {
				if msrp != price {
					percent := strings.TrimSpace(e.DOM.Find(".discount_pct").Text())
					price = price + " (" + percent + ")"
				}
			}
			if price != "" {
				data = append(data, TrendingGame{
					Title: title,
					Price: price,
					MSRP:  msrp,
				})
			}
		})
	})
	c.OnError(func(_ *colly.Response, err error) { check(err) })
	c.Visit("https://store.steampowered.com/search/?filter=topsellers")
	db.Exec("DELETE FROM trending_games")
	slog.Info(fmt.Sprintf("scraped %d games", len(data)))
	result := db.Create(data)
	slog.Info(fmt.Sprintf("inserted %d games", result.RowsAffected))
	check(result.Error)
}
