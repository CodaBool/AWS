package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func scrapeJS() {
	// log := logger.With().Str("func", "scrapeJS").Logger()

	defer wg.Done()
	var data []TrendingJS
	c := colly.NewCollector(colly.Async(true))
	c.OnHTML("main", func(e *colly.HTMLElement) {
		subjectSlice := strings.Split(fmt.Sprintf("%v", e.Request.URL), "%3A")
		pageSlice := strings.Split(fmt.Sprintf("%v", e.Request.URL), "page=")
		pageStr := strings.Split(pageSlice[1], "&q=")[0]
		page, err := strconv.Atoi(pageStr)
		check(err)
		e.ForEach("section", func(i int, el *colly.HTMLElement) {
			title := el.DOM.First().Find("h3").Text()
			description := el.DOM.First().Find("p").Text()
			data = append(data, TrendingJS{
				Title:       title,
				Description: description,
				Page:        page,
				Rank:        page*20 + i + 1,
				Subject:     subjectSlice[1],
			})
		})
	})

	c.OnError(func(_ *colly.Response, err error) { check(err) })
	// for page := 0; page < 2; page++ {
	for _, subject := range []string{"backend", "front-end", "cli", "framework"} {
		// log.Print("page ", page, ", ", subject)s
		log.Print(subject)
		// c.Visit("https://www.npmjs.com/search?ranking=popularity&page=" + strconv.Itoa(page) + "&q=keywords%3A" + subject)
		c.Visit("https://www.npmjs.com/search?ranking=popularity&page=0&q=keywords%3A" + subject)
	}
	// }
	c.Wait()
	db.Exec("DELETE FROM trending_js")
	slog.Info(fmt.Sprintf("+%d js", len(data)))
	db.Create(data)
}
