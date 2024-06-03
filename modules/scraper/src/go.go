package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

func scrapeLibhunt() []TrendingGo {
	// log := logger.With().Str("func", "scrapeLibhunt").Logger()

	c := colly.NewCollector(colly.Async(true))
	// c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
	var libData []TrendingGo
	c.OnHTML(".lib-list", func(e *colly.HTMLElement) {
		// fmt.Println("testing", e.Text)
		// e.ChildAttr("a", "href")
		// e.ChildAttr("img", "src")
		// e.ChildText("h2")
		// e.ChildText(".price")
		e.ForEach("li", func(i int, e *colly.HTMLElement) {
			// href := e.ChildAttr("a", "href")
			// splits := strings.Split(href, "/")
			slog.Debug(fmt.Sprintf("name %s", e.ChildText("h3")))
			libData = append(libData, TrendingGo{
				Name:        strings.ToLower(strings.TrimSpace(e.ChildText("h3"))),
				// FullName:    splits[3] + "/" + splits[4],
				// Href:        href,
				Description: strings.TrimSpace(e.ChildText(".tagline")),
			})
		})
	})

	c.OnError(func(_ *colly.Response, err error) {
		check(err)
	})

	for page := 1; page < 5; page++ {
		c.Visit("https://go.libhunt.com/projects?page=" + strconv.Itoa(page))
	}
	c.Wait()
	return libData
}

func scrapeGo() {
	// log := logger.With().Str("func", "scrapeGo").Logger()

	defer wg.Done()
	var ghData []TrendingGo
	client := &http.Client{Timeout: 9 * time.Second}
	for page := 1; page < 5; page++ {
		log.Println("scraping page ", page, " of go repos")
		time.Sleep(6 * time.Second)
		req, err := http.NewRequest(http.MethodGet, "https://api.github.com/search/repositories?q=language:golang&stars:%3E1&sort=stars&order=desc&per_page=100&page="+strconv.Itoa(page), nil)
		check(err)
		req.Header.Add("Authorization", os.Getenv("GIT_TOKEN"))
		res, err := client.Do(req)
		check(err)
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		check(err)

		gjson.GetBytes(body, "items").ForEach(func(_, item gjson.Result) bool {
			var repo TrendingGo
			item.ForEach(func(key, val gjson.Result) bool {
				if key.String() == "name" {
					repo.Name = strings.ToLower(val.String())
				} else if key.String() == "stargazers_count" {
					repo.Stars = val.Int()
				} else if key.String() == "full_name" {
					repo.FullName = val.String()
				}
				return true
			})
			ghData = append(ghData, repo)
			return true
		})
	}

	libData := scrapeLibhunt()

	matchCount := 0
	slog.Info("fetching stars for go repos (est. 4 minutes)")
	for keyLIB, valLIB := range libData {
		found := false
		for _, valGH := range ghData {
			if valGH.Name == valLIB.Name {
				log.Println(keyLIB, " match ", valLIB.Name)
				matchCount++
				found = true
				slog.Debug(fmt.Sprintf("%d/119 fetching stars", keyLIB))
				libData[keyLIB].Stars = valGH.Stars
			}
		}
		if !found {
			log.Println(keyLIB, " no match ", valLIB.Name)
			time.Sleep(6 * time.Second)
			req, err := http.NewRequest(http.MethodGet, "https://api.github.com/search/repositories?q="+valLIB.FullName, nil)
			check(err)
			req.Header.Add("Authorization", os.Getenv("GIT_TOKEN"))
			res, err := client.Do(req)
			check(err)
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			check(err)

			gjson.GetBytes(body, "items").ForEach(func(i, item gjson.Result) bool {
				item.ForEach(func(key, val gjson.Result) bool {
					if key.String() == "stargazers_count" {
						log.Println("  stars = ", val.Int())
						libData[keyLIB].Stars = val.Int()
					}
					return true
				})
				return false
			})
		}
	}
	log.Println("matched ", matchCount, "/", len(libData))

	db.Exec("DELETE FROM trending_gos")
	slog.Info(fmt.Sprintf("+%s go", strconv.Itoa(len(libData))))
	db.Create(libData)
}
