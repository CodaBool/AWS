package main

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/joho/godotenv/autoload"
)

var wg sync.WaitGroup

func main() {
	local := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == ""
	buildLogger(false, true, local)
	if local {
		handle(context.TODO(), nil)
	} else {
		lambda.Start(handle)
	}
}

func handle(ctx context.Context, _ any) (string, error) {
	dbInit(true)
	wg.Add(8)
	go scrapePY(false)        // +100 rows | false = 101 req | true = 1 req
	go scrapeGames()          // +49  rows | 1 req
	go scrapeGithub()         // +100 rows | 1 req
	go scrapeGo() // +120 rows | 30 req
	go scrapeUpcomingMovies() // +171 rows | 1 req
	go scrapeTV()             // +100 rows | 1 req
	go scrapeTrendingMovies() // +100 rows | 1 req
	go scrapeJS()             // +160 rows | 8 req
	wg.Wait()
	return "", nil
}
