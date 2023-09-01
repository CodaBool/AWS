package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func scrapeGithub() {
	// log := logger.With().Str("func", "scrapeGithub").Logger()

	defer wg.Done()
	var data []TrendingGithub
	client := &http.Client{Timeout: 9 * time.Second}

	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/search/repositories?q=stars:%3E1&sort=stars&order=desc&per_page=100", nil)
	check(err)
	req.Header.Add("Authorization", os.Getenv("GIT_TOKEN"))
	res, err := client.Do(req)
	check(err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	check(err)

	gjson.GetBytes(body, "items").ForEach(func(_, item gjson.Result) bool {
		var repo TrendingGithub
		item.ForEach(func(key, val gjson.Result) bool {
			if key.String() == "name" {
				repo.Name = strings.ToLower(val.String())
			} else if key.String() == "stargazers_count" {
				repo.Stars = val.Int()
			} else if key.String() == "description" {
				repo.Description = val.String()
			}
			return true
		})
		data = append(data, repo)
		return true
	})
	db.Exec("DELETE FROM trending_githubs")
	slog.Info(fmt.Sprintf("+%d github", len(data)))
	db.Create(data)
}
