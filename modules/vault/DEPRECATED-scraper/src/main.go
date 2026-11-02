package main

import (
	"context"
	"os"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	local := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == ""
	buildLogger(false, true, local)
	_ = os.MkdirAll("/tmp/.config", 0o755)
	_ = os.MkdirAll("/tmp/.cache", 0o755)
	_ = os.MkdirAll("/tmp/fontconfig", 0o755)
	if local {
		os.Setenv("CHROME_PATH", "/usr/bin/chromium-browser")
		handle(context.TODO(), nil)
	} else {
		lambda.Start(handle)
	}
}

func handle(ctx context.Context, _ any) (string, error) {
	dbInit(true)
	scrapePY(false)        // +100 rows | false = 101 req | true = 1 req
	scrapeGames()          // +49  rows | 1 req
	scrapeGithub()         // +100 rows | 1 req
	scrapeGo(false)             // +120 rows | 30 req | true = libhunt | false = no desc
	scrapeUpcomingMovies() // +171 rows | 1 req
	scrapeTV()             // +100 rows | 1 req
	scrapeTrendingMovies() // +100 rows | 1 req
	scrapeJS()             // +160 rows | 8 req
	return "", nil
}
