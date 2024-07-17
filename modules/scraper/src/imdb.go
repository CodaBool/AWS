package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"gorm.io/gorm"
)

func scrapeUpcomingMovies() {
	// log := logger.With().Str("func", "scrapeUpcomingMovies").Logger()

	defer wg.Done()
	var data []UpcomingMovie
	c := colly.NewCollector()
	c.OnHTML(".ipc-page-section--base article", func(e *colly.HTMLElement) {
		release := e.DOM.Children().Eq(0).Text()
		e.ForEach("ul", func(_ int, e *colly.HTMLElement) {
			e.ForEach("a", func(_ int, el *colly.HTMLElement) {
				if el.Text != "" {
					releaseTime, err := time.Parse("Jan 2, 2006", release)
					check(err)
					slog.Debug(fmt.Sprintf("Movies, get title from text %s", el.Text))
					data = append(data, UpcomingMovie{
						Release: releaseTime,
						Title:   el.Text[:len(el.Text)-7],
					})
				}
			})
		})
	})
	c.OnError(func(_ *colly.Response, err error) { check(err) })
	c.Visit("https://www.imdb.com/calendar/?region=US&type=MOVIE")
	db.Exec("DELETE FROM upcoming_movies")
	slog.Info(fmt.Sprintf("scraped %d upcoming movies", len(data)))
	result := db.Create(data)
	slog.Info(fmt.Sprintf("inserted %d upcoming movies", result.RowsAffected))
	check(result.Error)
}

func scrapeTV() {
	// log := logger.With().Str("func", "scrapeTV").Logger()

	defer wg.Done()
	var data []TrendingTV
	var data2 []TrendingTV
	c := colly.NewCollector()

	// IMDB is splitting traffic between 2 versions of the site
	// For now I am scraping both possible versions since
	// idk which one nginx will give me
	c.OnHTML(".compact-list-view", func(e *colly.HTMLElement) {
		e.ForEach("li", func(j int, el *colly.HTMLElement) {
			var tempData TrendingTV
			sign := "-"
			if el.DOM.Find(".rank-up").Length() == 1 {
				sign = "+"
			}
			vel := el.DOM.Find(".meter-const-ranking").Text()
			slog.Debug(fmt.Sprintf("TV, get velocity from %s", vel))
			vel = strings.Split(vel, "(")[1]
			vel = strings.Split(vel, ")")[0]
			slog.Debug(fmt.Sprintf("TV, final velocity %s", vel))
			if vel == "" {
				sign = ""
				vel = "no change"
			}
			tempData.Velocity = sign + vel
			tempData.Rank = j + 1
			tempData.Title = strings.TrimSpace(el.DOM.Find("h3").Text())
			tempData.Rating = strings.TrimSpace(el.DOM.Find(".ipc-rating-star--imdb").Text())
			data = append(data, tempData)
		})
	})
	c.OnHTML(".lister-list", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, e *colly.HTMLElement) {
			e.ForEach("a", func(_ int, el *colly.HTMLElement) {
				if strings.TrimSpace(el.Text) != "" {
					sign := "+"
					if e.DOM.Find(".down").Length() == 1 {
						sign = "-"
					}
					vel := strings.ReplaceAll(e.DOM.Find(".velocity").Text(), "\n", "")
					slog.Debug(fmt.Sprintf("TV trend, get velocity from %s", vel))
					vel = strings.Split(vel, "(")[1]
					if vel == "no change)" {
						sign = ""
						vel = "0)"
					}
					slog.Debug(fmt.Sprintf("TV trend, get velocity from %s at %d", vel, len(vel)-1))
					data2 = append(data2, TrendingTV{
						Title:    el.Text,
						Rank:     i,
						Velocity: sign + vel[:len(vel)-1],
						Rating:   strings.TrimSpace(e.DOM.Find(".imdbRating").Text()),
					})
				}
			})
		})
	})
	c.OnError(func(_ *colly.Response, err error) { check(err) })
	c.Visit("https://www.imdb.com/chart/tvmeter")
	db.Exec("DELETE FROM trending_tvs")
	slog.Info(fmt.Sprintf("new site scraped = %d | old site scraped = %d", len(data), len(data2)))
	// slog.Info(fmt.Sprintf("+%s tv", strconv.Itoa(len(data)+len(data2))))
	var result *gorm.DB
	if len(data) > len(data2) {
		result = db.Create(data)
	} else {
		result = db.Create(data2)
	}
	slog.Info(fmt.Sprintf("inserted %d tv", result.RowsAffected))
	check(result.Error)
}

func scrapeTrendingMovies() {
	// log := logger.With().Str("func", "scrapeTrendingMovies").Logger()
	defer wg.Done()
	var data []TrendingMovie
	var data2 []TrendingMovie
	c := colly.NewCollector()

	// IMDB is splitting traffic between 2 versions of the site
	// For now I am scraping both possible versions since
	// idk which one nginx will give me
	c.OnHTML(".compact-list-view", func(e *colly.HTMLElement) {
		e.ForEach("li", func(j int, el *colly.HTMLElement) {
			var tempData TrendingMovie
			sign := "-"
			if el.DOM.Find(".rank-up").Length() == 1 {
				sign = "+"
			}
			vel := el.DOM.Find(".meter-const-ranking").Text()
			slog.Debug(fmt.Sprintf("movies trend, get velocity from %s", vel))
			vel = strings.Split(vel, "(")[1]
			vel = strings.Split(vel, ")")[0]
			slog.Debug(fmt.Sprintf("movies trend, final velocity %s", vel))
			if vel == "" {
				sign = ""
				vel = "no change"
			}

			tempData.Velocity = sign + vel
			tempData.Rank = j + 1
			tempData.Title = strings.TrimSpace(el.DOM.Find("h3").Text())
			tempData.Rating = strings.TrimSpace(el.DOM.Find(".ipc-rating-star--imdb").Text())
			data = append(data, tempData)
		})
	})
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
					slog.Debug(fmt.Sprintf("POP MOVIES, get velocity from %s", vel))
					vel = strings.Split(vel, "(")[1]
					if vel == "no change)" {
						sign = ""
						vel = "0)"
					}
					slog.Debug(fmt.Sprintf("POP MOVIES, get velocity from %s at %d", vel, len(vel)-1))
					tempData.Velocity = sign + vel[:len(vel)-1]
					tempData.Rank = j + 1
					tempData.Title = ele.DOM.Children().First().Text()
				} else if i == 2 {
					tempData.Rating = strings.TrimSpace(ele.Text)
				}
			})
			data2 = append(data2, tempData)
		})
	})
	c.OnError(func(_ *colly.Response, err error) { check(err) })
	c.Visit("https://www.imdb.com/chart/moviemeter")
	db.Exec("DELETE FROM trending_movies")

	// again a split on site version means idk which data slice will have data
	slog.Info(fmt.Sprintf("new site scraped = %d | old site scraped = %d", len(data), len(data2)))
	slog.Info(fmt.Sprintf("scraped %d trending movies", len(data)+len(data2)))

	var result *gorm.DB
	if len(data) > len(data2) {
		result = db.Create(data)
	} else {
		result = db.Create(data2)
	}
	slog.Info(fmt.Sprintf("inserted %d trending movies", result.RowsAffected))
	check(result.Error)
}
