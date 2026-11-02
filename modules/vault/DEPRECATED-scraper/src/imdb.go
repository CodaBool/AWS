package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	// "github.com/gocolly/colly"
	// "gorm.io/gorm"
)

func scrapeUpcomingMovies() {
	var data []UpcomingMovie

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

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var html string

	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate("https://www.imdb.com/calendar/?region=US&type=MOVIE"),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.WaitVisible(`article[data-testid="calendar-section"]`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
		chromedp.Evaluate(`document.documentElement.outerHTML`, &html),
	)
	if err != nil {
		slog.Error("chromedp failed", "err", err)
		return
	}

	slog.Info("captured html", "len", len(html))

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		slog.Error("goquery parse failed", "err", err)
		return
	}

	doc.Find(`article[data-testid="calendar-section"]`).Each(func(_ int, sec *goquery.Selection) {
		releaseText := strings.TrimSpace(
			sec.Find(`[data-testid="release-date"] .ipc-title__text`).First().Text(),
		)
		if releaseText == "" {
			return
		}

		releaseTime, err := time.Parse("Jan 2, 2006", releaseText)
		if err != nil {
			slog.Error("failed to parse release date", "releaseText", releaseText, "err", err)
			return
		}

		sec.Find(`li[data-testid="coming-soon-entry"]`).Each(func(_ int, item *goquery.Selection) {
			title := strings.TrimSpace(
				item.Find(`a.ipc-metadata-list-summary-item__t`).First().Text(),
			)
			if title == "" {
				title, _ = item.Find("a").First().Attr("aria-label")
				title = strings.TrimSpace(title)
			}
			if title == "" {
				return
			}

			data = append(data, UpcomingMovie{
				Title:   title,
				Release: releaseTime,
			})
		})
	})

	slog.Info("scraped upcoming movies", "count", len(data))

	if len(data) == 0 {
		slog.Warn("no upcoming movies scraped")
		return
	}

	if err := db.Exec("DELETE FROM upcoming_movies").Error; err != nil {
		slog.Error("delete failed", "err", err)
		return
	}

	result := db.Create(&data)
	if result.Error != nil {
		slog.Error("insert failed", "err", result.Error)
		return
	}

	slog.Info("inserted upcoming movies", "rows", result.RowsAffected)
}

func scrapeTV() {
	// defer wg.Done()

	var data []TrendingTV
	var data2 []TrendingTV

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

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var html string

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return network.Enable().Do(ctx)
		}),
		chromedp.Navigate("https://www.imdb.com/chart/tvmeter"),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)
	if err != nil {
		slog.Error("chromedp failed", "err", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		slog.Error("goquery parse failed", "err", err)
		return
	}

	// Newer site version
	doc.Find(".compact-list-view").Each(func(_ int, section *goquery.Selection) {
		section.Find("li").Each(func(j int, item *goquery.Selection) {
			var tempData TrendingTV

			sign := "-"
			if item.Find(".rank-up").Length() == 1 {
				sign = "+"
			}

			vel := strings.TrimSpace(item.Find(".meter-const-ranking").Text())
			slog.Debug(fmt.Sprintf("TV, get velocity from %s", vel))

			if start := strings.Index(vel, "("); start >= 0 {
				vel = vel[start+1:]
				if end := strings.Index(vel, ")"); end >= 0 {
					vel = vel[:end]
				}
			} else {
				vel = ""
			}

			slog.Debug(fmt.Sprintf("TV, final velocity %s", vel))

			if vel == "" {
				sign = ""
				vel = "no change"
			}

			tempData.Velocity = sign + vel
			tempData.Rank = j + 1
			tempData.Title = strings.TrimSpace(item.Find("h3").First().Text())
			tempData.Rating = strings.TrimSpace(item.Find(".ipc-rating-star--imdb").First().Text())

			if tempData.Title != "" {
				data = append(data, tempData)
			}
		})
	})

	// Older site version
	doc.Find(".lister-list").Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			title := strings.TrimSpace(row.Find("a").First().Text())
			if title == "" {
				return
			}

			sign := "+"
			if row.Find(".down").Length() == 1 {
				sign = "-"
			}

			vel := strings.ReplaceAll(strings.TrimSpace(row.Find(".velocity").Text()), "\n", "")
			slog.Debug(fmt.Sprintf("TV trend, get velocity from %s", vel))

			if start := strings.Index(vel, "("); start >= 0 {
				vel = vel[start+1:]
			} else {
				vel = ""
			}

			if vel == "no change)" {
				sign = ""
				vel = "0)"
			}

			if strings.HasSuffix(vel, ")") {
				vel = vel[:len(vel)-1]
			}

			slog.Debug(fmt.Sprintf("TV trend, final velocity %s", vel))

			data2 = append(data2, TrendingTV{
				Title:    title,
				Rank:     i,
				Velocity: sign + vel,
				Rating:   strings.TrimSpace(row.Find(".imdbRating").First().Text()),
			})
		})
	})


	var chosen []TrendingTV
	if len(data) > len(data2) {
		chosen = data
	} else {
		chosen = data2
	}

	if len(chosen) == 0 {
		slog.Warn("no trending tv scraped")
		return
	}

	db.Exec("DELETE FROM trending_tvs")
	slog.Info(fmt.Sprintf("new site scraped = %d | old site scraped = %d", len(data), len(data2)))

	result := db.Create(&chosen)
	if result.Error != nil {
		slog.Error("insert failed", "err", result.Error)
		return
	}

	slog.Info(fmt.Sprintf("inserted %d tv", result.RowsAffected))
}

func scrapeTrendingMovies() {
	// defer wg.Done()

	var data []TrendingMovie
	var data2 []TrendingMovie

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

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var html string

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return network.Enable().Do(ctx)
		}),
		chromedp.Navigate("https://www.imdb.com/chart/moviemeter"),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)
	if err != nil {
		slog.Error("chromedp failed", "err", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		slog.Error("goquery parse failed", "err", err)
		return
	}

	// Newer site version
	doc.Find(".compact-list-view").Each(func(_ int, section *goquery.Selection) {
		section.Find("li").Each(func(j int, item *goquery.Selection) {
			var tempData TrendingMovie

			sign := "-"
			if item.Find(".rank-up").Length() == 1 {
				sign = "+"
			}

			vel := strings.TrimSpace(item.Find(".meter-const-ranking").Text())
			slog.Debug(fmt.Sprintf("movies trend, get velocity from %s", vel))

			if start := strings.Index(vel, "("); start >= 0 {
				vel = vel[start+1:]
				if end := strings.Index(vel, ")"); end >= 0 {
					vel = vel[:end]
				}
			} else {
				vel = ""
			}

			slog.Debug(fmt.Sprintf("movies trend, final velocity %s", vel))

			if vel == "" {
				sign = ""
				vel = "no change"
			}

			tempData.Velocity = sign + vel
			tempData.Rank = j + 1
			tempData.Title = strings.TrimSpace(item.Find("h3").First().Text())
			tempData.Rating = strings.TrimSpace(item.Find(".ipc-rating-star--imdb").First().Text())

			if tempData.Title != "" {
				data = append(data, tempData)
			}
		})
	})

	// Older site version
	doc.Find(".lister-list").Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			var tempData TrendingMovie

			row.Find("td").Each(func(i int, cell *goquery.Selection) {
				if i == 1 {
					sign := "+"
					if cell.Find(".down").Length() == 1 {
						sign = "-"
					}

					vel := strings.ReplaceAll(strings.TrimSpace(cell.Find(".velocity").Text()), "\n", "")
					slog.Debug(fmt.Sprintf("POP MOVIES, get velocity from %s", vel))

					if start := strings.Index(vel, "("); start >= 0 {
						vel = vel[start+1:]
					} else {
						vel = ""
					}

					if vel == "no change)" {
						sign = ""
						vel = "0)"
					}

					if strings.HasSuffix(vel, ")") {
						vel = vel[:len(vel)-1]
					}

					slog.Debug(fmt.Sprintf("POP MOVIES, final velocity %s", vel))

					tempData.Velocity = sign + vel
					tempData.Rank = j + 1
					tempData.Title = strings.TrimSpace(cell.Children().First().Text())
				} else if i == 2 {
					tempData.Rating = strings.TrimSpace(cell.Text())
				}
			})

			if tempData.Title != "" {
				data2 = append(data2, tempData)
			}
		})
	})

	db.Exec("DELETE FROM trending_movies")

	slog.Info(fmt.Sprintf("new site scraped = %d | old site scraped = %d", len(data), len(data2)))
	slog.Info(fmt.Sprintf("scraped %d trending movies", len(data)+len(data2)))

	var chosen []TrendingMovie
	if len(data) > len(data2) {
		chosen = data
	} else {
		chosen = data2
	}

	if len(chosen) == 0 {
		slog.Warn("no trending movies scraped")
		return
	}

	result := db.Create(&chosen)
	if result.Error != nil {
		slog.Error("insert failed", "err", result.Error)
		return
	}

	slog.Info(fmt.Sprintf("inserted %d trending movies", result.RowsAffected))
}
