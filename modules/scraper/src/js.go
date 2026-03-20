package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func scrapeJS() {
	// defer wg.Done()

	var data []TrendingJS

	chromePath := os.Getenv("CHROME_PATH")
	if chromePath == "" {
		chromePath = "/usr/bin/chromium"
	}

	profileDir, errProf := os.MkdirTemp("/tmp", "chromedp-profile-*")
	if errProf != nil {
		slog.Error("chrome profile", "err", errProf)
		return
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
		chromedp.Flag("disable-background-networking", true),
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

	subjects := []string{"backend", "front-end", "cli", "framework"}

	for _, subject := range subjects {
		page := 0
		url := "https://npmjs.com/search?ranking=popularity&page=" + strconv.Itoa(page) + "&q=keywords%3A" + url.QueryEscape(subject)

		slog.Debug(fmt.Sprintf("section %s", subject))

		var html string
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(2*time.Second),
			chromedp.OuterHTML("html", &html, chromedp.ByQuery),
		)
		if err != nil {
			slog.Error("chromedp failed", "subject", subject, "err", err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			slog.Error("goquery parse failed", "subject", subject, "err", err)
			continue
		}

		// npm search results are rendered as package cards; each one generally has an h3 and a p.
		// We keep this a little defensive so minor markup shifts do not kill the scrape.
		rankBase := page * 20

		doc.Find("section").Each(func(i int, sec *goquery.Selection) {
			title := strings.TrimSpace(sec.Find("h3").First().Text())
			description := strings.TrimSpace(sec.Find("p").First().Text())

			if title == "" {
				return
			}

			data = append(data, TrendingJS{
				Title:       title,
				Description: description,
				Page:        page,
				Rank:        rankBase + i + 1,
				Subject:     subject,
			})
		})
	}

	db.Exec("DELETE FROM trending_js")
	slog.Info(fmt.Sprintf("scraped %d js", len(data)))

	if len(data) == 0 {
		slog.Warn("no trending js scraped")
		return
	}

	result := db.Create(&data)
	if result.Error != nil {
		slog.Error("insert failed", "err", result.Error)
		return
	}

	slog.Info(fmt.Sprintf("inserted %d js", result.RowsAffected))
}
