package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

var dg *discordgo.Session

var channel = "1254921386267250879"

type Input struct {
	Test bool `json:"test"`
}

func main() {
	local := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == ""
	buildLogger(true, false, local)
	if local {
		handle(context.TODO(), Input{
			Test: true,
		})
	} else {
		lambda.Start(handle)
	}
}

func handle(ctx context.Context, i Input) (string, error) {
	if i.Test {
		channel = "870190331554054194"
	}

	now := time.Now()
	year, month, _ := now.Date()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	firstSaturdayOffset := (6 - int(firstDay.Weekday()) + 7) % 7
	secondSaturday := firstDay.AddDate(0, 0, firstSaturdayOffset+7)

	// UTC -> 2:30 PM EST
	secondSaturday = secondSaturday.Add(time.Hour*19 + time.Minute*30)
	slog.Info("target time is " + secondSaturday.Add(-(time.Hour * 5)).Format(time.RFC3339) + " EST")
	slog.Info(fmt.Sprintf("today is %dth day at hour %d", now.Day(), now.Hour()-5))

	threeDaysBefore := secondSaturday.AddDate(0, 0, -3)
	oneDayBefore := secondSaturday.AddDate(0, 0, -1)
	oneHourBefore := secondSaturday.Add(-time.Hour)

	timestampFull := fmt.Sprintf("<t:%d:f>", secondSaturday.Unix())
	timestampRel := fmt.Sprintf("<t:%d:R>", secondSaturday.Unix())
	bot, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	check(err)

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
