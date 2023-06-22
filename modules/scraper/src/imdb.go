package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func scrapeUpcomingMovies() {
	log := logger.With().Str("func", "scrapeUpcomingMovies").Logger()

	defer wg.Done()
	var data []UpcomingMovie
	c := colly.NewCollector()
	c.OnHTML(".ipc-page-section--base article", func(e *colly.HTMLElement) {
		release := e.DOM.Children().Eq(0).Text()
		e.ForEach("ul", func(_ int, e *colly.HTMLElement) {
			e.ForEach("a", func(_ int, el *colly.HTMLElement) {
				if el.Text != "" {
					releaseTime, err := time.Parse("Jan 2, 2006", release)
					check(err, log)
					data = append(data, UpcomingMovie{
						Release: releaseTime,
						Title:   el.Text[:len(el.Text)-7],
					})
				}
			})
		})
	})
	c.OnError(func(_ *colly.Response, err error) { check(err, log) })
	c.Visit("https://www.imdb.com/calendar/?region=US&type=MOVIE")
	db.Exec("DELETE FROM upcoming_movies")
	log.Info().Msg(fmt.Sprintf("+%d upcoming movies", len(data)))
	db.Create(data)
}

func scrapeTV() {
	log := logger.With().Str("func", "scrapeTV").Logger()

	defer wg.Done()
	var data []TrendingTV
	c := colly.NewCollector()
	c.OnHTML(".lister-list", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, e *colly.HTMLElement) {
			e.ForEach("a", func(_ int, el *colly.HTMLElement) {
				if strings.TrimSpace(el.Text) != "" {
					sign := "+"
					if e.DOM.Find(".down").Length() == 1 {
						sign = "-"
					}
					vel := strings.ReplaceAll(e.DOM.Find(".velocity").Text(), "\n", "")
					vel = strings.Split(vel, "(")[1]
					if vel == "no change)" {
						sign = ""
						vel = "0)"
					}
					data = append(data, TrendingTV{
						Title:    el.Text,
						Rank:     i,
						Velocity: sign + vel[:len(vel)-1],
						Rating:   strings.TrimSpace(e.DOM.Find(".imdbRating").Text()),
					})
				}
			})
		})
	})
	c.OnError(func(_ *colly.Response, err error) { check(err, log) })
	c.Visit("https://www.imdb.com/chart/tvmeter")
	db.Exec("DELETE FROM trending_tvs")
	log.Info().Msg(fmt.Sprintf("+%d tv", len(data)))
	db.Create(data)
}

func scrapeTrendingMovies() {
	log := logger.With().Str("func", "scrapeTrendingMovies").Logger()
	defer wg.Done()
	var data []TrendingMovie
	c := colly.NewCollector()
	c.OnHTML(".lister-list", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(j int, el *colly.HTMLElement) {
			var tempData TrendingMovie
			el.ForEach("td", func(i int, ele *colly.HTMLElement) {
				if i == 1 {
					sign := "+"
					if ele.DOM.Find(".down").Length() == 1 {
						sign = "-"
					}
					vel := strings.ReplaceAll(ele.DOM.Find(".velocity").Text(), "\n", "")
					vel = strings.Split(vel, "(")[1]
					if vel == "no change)" {
						sign = ""
						vel = "0)"
					}
					tempData.Velocity = sign + vel[:len(vel)-1]
					tempData.Rank = j + 1
					tempData.Title = ele.DOM.Children().First().Text()
				} else if i == 2 {
					tempData.Rating = strings.TrimSpace(ele.Text)
				}
			})
			data = append(data, tempData)
		})
	})
	c.OnError(func(_ *colly.Response, err error) { check(err, log) })
	c.Visit("https://www.imdb.com/chart/moviemeter")
	db.Exec("DELETE FROM trending_movies")
	log.Info().Msg(fmt.Sprintf("+%d trending movies", len(data)))
	db.Create(data)
}
