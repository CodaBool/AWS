package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
)

func scrapeLibhunt() []TrendingGo {
	var libData []TrendingGo

	chromePath := os.Getenv("CHROME_PATH")
	if chromePath == "" {
		chromePath = "/usr/bin/chromium"
	}

	profileDir, errProf := os.MkdirTemp("/tmp", "chromedp-profile-*")
	if errProf != nil {
		slog.Error("chrome profile", "err", errProf)
		return nil
	}
	defer os.RemoveAll(profileDir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.Flag("headless", "new"),
		chromedp.NoSandbox,
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-zygote", true),
		chromedp.Flag("single-process", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("disable-crash-reporter", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("user-data-dir", profileDir),
		chromedp.Flag("data-path", profileDir),
		chromedp.Flag("disk-cache-dir", "/tmp/chromedp-cache"),
		chromedp.Flag("window-size", "1366,768"),
		chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	for page := 1; page < 5; page++ {
		url := "https://go.libhunt.com/projects?page=" + strconv.Itoa(page)
		var html string

		err := chromedp.Run(ctx,
			network.Enable(),
			chromedp.Navigate(url),
			chromedp.WaitReady("body", chromedp.ByQuery),
			chromedp.Sleep(2*time.Second),
			chromedp.Evaluate(`document.documentElement.outerHTML`, &html),
		)
		if err != nil {
			slog.Error("chromedp libhunt failed", "page", page, "err", err)
			continue
		}

		if strings.Contains(strings.ToLower(html), "just a moment") {
			slog.Error("cloudflare challenge hit in browser too", "page", page, "url", url)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			slog.Error("goquery parse failed", "page", page, "err", err)
			continue
		}

		doc.Find(".lib-list li").Each(func(_ int, el *goquery.Selection) {
			rawName := strings.TrimSpace(el.Find("h3").First().Text())
			if rawName == "" {
				return
			}

			name := strings.ReplaceAll(strings.ToLower(rawName), " ", "-")
			fullName := name

			if href, ok := el.Find(`a[target="_blank"]`).First().Attr("href"); ok {
				splits := strings.Split(strings.TrimSpace(href), "/")
				if len(splits) >= 5 {
					fullName = splits[3] + "/" + splits[4]
				}
			}

			libData = append(libData, TrendingGo{
				Name:        name,
				FullName:    fullName,
				Description: strings.TrimSpace(el.Find(".tagline").First().Text()),
			})
		})

		time.Sleep(1500 * time.Millisecond)
	}

	return libData
}

func scrapeGo(includeLibhunt bool) {
	var ghData []TrendingGo
	client := &http.Client{Timeout: 9 * time.Second}

	for page := 1; page < 5; page++ {
		slog.Info(fmt.Sprintf("scraping page %d of go repos", page))
		time.Sleep(3 * time.Second)

		req, err := http.NewRequest(
			http.MethodGet,
			"https://api.github.com/search/repositories?q=language:golang&stars:%3E1&sort=stars&order=desc&per_page=100&page="+strconv.Itoa(page),
			nil,
		)
		check(err)

		req.Header.Add("Authorization", os.Getenv("GIT_TOKEN"))
		req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
		req.Header.Add("Accept", "application/vnd.github+json")

		res, err := client.Do(req)
		check(err)
		slog.Info(fmt.Sprintf("HTTP code: %d", res.StatusCode))

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		check(err)

		gjson.GetBytes(body, "items").ForEach(func(_, item gjson.Result) bool {
			var repo TrendingGo
			item.ForEach(func(key, val gjson.Result) bool {
				switch key.String() {
				case "name":
					repo.Name = strings.ToLower(val.String())
				case "stargazers_count":
					repo.Stars = val.Int()
				case "full_name":
					repo.FullName = val.String()
				case "description":
					repo.Description = val.String()
				}
				return true
			})
			ghData = append(ghData, repo)
			return true
		})
	}

	var finalData []TrendingGo

	// libhunt uses cloudflare bot detection, make optional
	if includeLibhunt {
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
					libData[keyLIB].Stars = valGH.Stars

					// optionally also take GitHub description if LibHunt one is empty
					if strings.TrimSpace(libData[keyLIB].Description) == "" {
						libData[keyLIB].Description = valGH.Description
					}
					break
				}
			}

			if !found {
				slog.Info(fmt.Sprintf("%d no match %s, using github api for stars", keyLIB, valLIB.Name))
				time.Sleep(15 * time.Second)

				q := valLIB.Name
				if valLIB.FullName != "" {
					q = valLIB.FullName
				}

				req, err := http.NewRequest(
					http.MethodGet,
					"https://api.github.com/search/repositories?q="+q,
					nil,
				)
				check(err)

				req.Header.Add("Authorization", os.Getenv("GIT_TOKEN"))
				req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
				req.Header.Add("Accept", "application/vnd.github+json")

				res, err := client.Do(req)
				check(err)
				slog.Info(fmt.Sprintf("HTTP code: %d", res.StatusCode))

				body, err := io.ReadAll(res.Body)
				res.Body.Close()
				check(err)

				if res.StatusCode == 403 {
					slog.Error("bad token or rate limited")
				} else if res.StatusCode == 422 {
					slog.Error("bad query or validation failed")
				}

				gjson.GetBytes(body, "items").ForEach(func(i, item gjson.Result) bool {
					item.ForEach(func(key, val gjson.Result) bool {
						switch key.String() {
						case "stargazers_count":
							slog.Info(fmt.Sprintf("stars = %d", val.Int()))
							libData[keyLIB].Stars = val.Int()
						case "description":
							if strings.TrimSpace(libData[keyLIB].Description) == "" {
								libData[keyLIB].Description = val.String()
							}
						}
						return true
					})
					return false
				})
			}
		}

		slog.Info(fmt.Sprintf("matched %d/%d", matchCount, len(libData)))
		finalData = libData
	} else {
		slog.Info("libhunt disabled, using github results only")
		finalData = ghData
	}

	if err := db.Exec("DELETE FROM trending_gos").Error; err != nil {
		slog.Error("failed clearing trending_gos", "err", err)
		return
	}

	slog.Info(fmt.Sprintf("attempting insert of %d go rows", len(finalData)))

	result := db.Create(&finalData)
	slog.Info(fmt.Sprintf("rows: %d", result.RowsAffected))
	check(result.Error)
}
