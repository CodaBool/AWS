package main

import (
	"fmt"
	"io"
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
			name := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(e.ChildText("h3"))), " ", "-")
			fullName := name
			if e.ChildAttr("a[target='_blank']", "href") != "" {
				// slog.Info(fmt.Sprintf("anchor with target=_blank found: %s", e.ChildAttr("a[target='_blank']", "href")))

				href := e.ChildAttr("a[target='_blank']", "href")
				splits := strings.Split(href, "/")
				// slog.Info(fmt.Sprintf("splits: %v", splits))
				// slog.Info(fmt.Sprintf("repo: %s/%s", splits[3], splits[4]))
				fullName = splits[3] + "/" + splits[4]
			}
			libData = append(libData, TrendingGo{
				Name:     name,
				FullName: fullName,
				// FullName:    splits[3] + "/" + splits[4],
				// Href:        href,
				Description: strings.TrimSpace(e.ChildText(".tagline")),
			})
		})
	})

	c.OnError(func(_ *colly.Response, err error) {
		slog.Error(fmt.Sprintf("go error: %s", err.Error()))
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
		slog.Info(fmt.Sprintf("scraping page %d of go repos", page))
		time.Sleep(4 * time.Second)
		req, err := http.NewRequest(http.MethodGet, "https://api.github.com/search/repositories?q=language:golang&stars:%3E1&sort=stars&order=desc&per_page=100&page="+strconv.Itoa(page), nil)
		check(err)
		req.Header.Add("Authorization", os.Getenv("GIT_TOKEN"))
		req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
		req.Header.Add("Accept", "application/vnd.github+json")
		res, err := client.Do(req)
		slog.Info(fmt.Sprintf("HTTP code: %d", res.StatusCode))
		check(err)
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
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
				slog.Info(fmt.Sprintf("%d match %s", keyLIB, valLIB.Name))
				matchCount++
				found = true
				slog.Debug(fmt.Sprintf("%d/119 fetching stars", keyLIB))
				libData[keyLIB].Stars = valGH.Stars
			}
		}
		if !found {
			slog.Info(fmt.Sprintf("%d no match %s, using github api for stars", keyLIB, valLIB.Name))
			// log.Println(keyLIB, " no match ", valLIB.Name)
			time.Sleep(15 * time.Second)
			q := valLIB.Name
			if valLIB.FullName != "" {
				q = valLIB.FullName
			}
			q = strings.ReplaceAll(q, " ", " ")
			req, err := http.NewRequest(http.MethodGet, "https://api.github.com/search/repositories?q="+q, nil)
			check(err)
			req.Header.Add("Authorization", os.Getenv("GIT_TOKEN"))
			req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
			req.Header.Add("Accept", "application/vnd.github+json")
			res, err := client.Do(req)
			slog.Info(fmt.Sprintf("HTTP code: %d", res.StatusCode))
			check(err)
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			check(err)

			if res.StatusCode == 422 {
				slog.Error("likely rate limited")
			} else if res.StatusCode == 422 {
				slog.Error("bad token or rate limited")
			}

			// slog.Info(fmt.Sprintf("Response Body: %s", string(body)))

			gjson.GetBytes(body, "items").ForEach(func(i, item gjson.Result) bool {
				item.ForEach(func(key, val gjson.Result) bool {
					if key.String() == "stargazers_count" {
						// log.Println("  stars = ", val.Int())
						slog.Info(fmt.Sprintf("  stars = %d", val.Int()))

						libData[keyLIB].Stars = val.Int()
					}
					return true
				})
				return false
			})
		}
	}
	slog.Info(fmt.Sprintf("matched %d/%d", matchCount, len(libData)))

	db.Exec("DELETE FROM trending_gos")
	slog.Info(fmt.Sprintf("attempting insert of +%s go rows", strconv.Itoa(len(libData))))

	result := db.Create(libData) // pass pointer of data to Create
	slog.Info(fmt.Sprintf("rows: %d", result.RowsAffected))
	check(result.Error)
}
