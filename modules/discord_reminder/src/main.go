package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

var dg *discordgo.Session

var channel = "1254921386267250879"

func main() {
	local := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == ""
	buildLogger(true, false, local)
	if local {
		handle(context.TODO(), events.LambdaFunctionURLRequest{
			QueryStringParameters: map[string]string{
				"body":   "wow",
				"action": "manual",
				"test":   "true",
			},
		})
	} else {
		lambda.Start(handle)
	}
}

func handle(ctx context.Context, req events.LambdaFunctionURLRequest) (string, error) {
	queryParams := req.QueryStringParameters
	action := queryParams["action"]
	secret := queryParams["secret"]
	test := queryParams["test"]
	body := queryParams["body"]

	if action == "" {
		return "", nil
	}
	if action == "manual" && secret != os.Getenv("TOKEN") {
		return "unauthorized", nil
	}

	if test == "true" {
		channel = "870190331554054194"
	}

	bot, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	check(err)

	if action == "manual" {
		_, err3 := bot.ChannelMessageSend(channel, body)
		check(err3)
		return "message sent", nil
	}

	now := time.Now()
	year, month, _ := now.Date()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	firstSaturdayOffset := (6 - int(firstDay.Weekday()) + 7) % 7
	secondSaturday := firstDay.AddDate(0, 0, firstSaturdayOffset+7)

	// UTC -> 2:30 PM EST
	loc, err := time.LoadLocation("America/New_York")
	check(err)
	secondSaturday = time.Date(secondSaturday.Year(), secondSaturday.Month(), secondSaturday.Day(), 14, 30, 0, 0, loc)
	slog.Info("target time is " + secondSaturday.Format(time.RFC3339) + " local")

	slog.Info(fmt.Sprintf("today is %dth day at hour %d", now.Day(), now.Hour()-5))

	threeDaysBefore := secondSaturday.AddDate(0, 0, -3)
	oneDayBefore := secondSaturday.AddDate(0, 0, -1)
	oneHourBefore := secondSaturday.Add(-time.Hour)

	timestampFull := fmt.Sprintf("<t:%d:f>", secondSaturday.Unix())
	timestampRel := fmt.Sprintf("<t:%d:R>", secondSaturday.Unix())

	if now.Year() == threeDaysBefore.Year() && now.Month() == threeDaysBefore.Month() && now.Day() == threeDaysBefore.Day() {
		slog.Info("3 days")
		_, err2 := bot.ChannelMessageSend(channel, "@everyone next session in "+timestampRel+","+timestampFull)
		check(err2)
	} else if now.Year() == oneDayBefore.Year() && now.Month() == oneDayBefore.Month() && now.Day() == oneDayBefore.Day() {
		slog.Info("1 day")
		_, err2 := bot.ChannelMessageSend(channel, "@everyone session "+timestampRel)
		check(err2)
	} else if now.Year() == secondSaturday.Year() && now.Month() == secondSaturday.Month() && now.Day() == secondSaturday.Day() && now.Hour() == oneHourBefore.Hour() {
		slog.Info("1 hour")
		_, err2 := bot.ChannelMessageSend(channel, "@everyone session in "+timestampRel)
		check(err2)
	} else {
		slog.Info("no reminder today")
	}

	return "", nil
}
