package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func scrapeGames() {
	defer wg.Done()
	var data []TrendingGame
	c := colly.NewCollector()
	c.OnHTML("#search_resultsRows", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, e *colly.HTMLElement) {
			title := e.DOM.Find(".search_name").Text()
			title = strings.ReplaceAll(strings.TrimSpace(title), "\n", "")
			title = strings.TrimSpace(strings.ReplaceAll(title, "VR Supported", ""))
			title = nonAlphanumericRegex.ReplaceAllString(title, "")
			price := e.DOM.Find(".search_price").Text()
			price = strings.TrimSpace(price)
			if price != "" {
				priceSlice := strings.Split(price, "$")
				if len(priceSlice) == 3 {
					percent := strings.TrimSpace(e.DOM.Find(".search_discount").Text())
					price = "$" + priceSlice[2] + " (" + percent + " off)"
				}
				log.Print(priceSlice)
				if strings.ToLower(priceSlice[0]) == "free to play" {
					price = "Free"
				}
				data = append(data, TrendingGame{
					Title: title,
					Price: price,
				})
			}
		})
	})
	c.OnError(func(_ *colly.Response, err error) { check(err) })
	c.Visit("https://store.steampowered.com/search/?filter=topsellers")
	db.Exec("DELETE FROM trending_games")
	log.Info().Msg(fmt.Sprintf("+%d games", len(data)))
	db.Create(data)
}
