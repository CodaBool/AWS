package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

var dg *discordgo.Session

var channel = "1254921386267250879"

var imageUrls = []string{
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_2.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_3.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_4.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_5.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_6.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_7.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_8.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_9.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_10.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_11.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_12.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_13.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/hq_14.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/sq_5.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/sq_6.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/sq_8.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/sq_11.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/sq_13.gif?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/sq_14.webp?raw=true",
	"https://github.com/CodaBool/cloudflare/blob/main/cron/img/sq_15.gif?raw=true",
}

func main() {
	local := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == ""
	buildLogger(true, false, local)
	if local {
		handle(context.TODO(), events.LambdaFunctionURLRequest{
			QueryStringParameters: map[string]string{
				"body":   "wow",              // any
				"action": "other",            // "manual", !=
				"test":   "true",             // "true", !=
				"secret": os.Getenv("TOKEN"), // match with what's in .env
			},
		})
	} else {
		lambda.Start(handle)
	}
}

func handle(ctx context.Context, req events.LambdaFunctionURLRequest) (string, error) {
	sess := session.Must(session.NewSession())
	ssmClient := ssm.New(sess)

	param, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String("post_discord_reminder"),
		WithDecryption: aws.Bool(true),
	})
	check(err)

	if *param.Parameter.Value != "true" {
		slog.Info("Reminders are paused. Exiting.")
		return "", nil
	} else {
		slog.Info("post_discord_reminder = " + *param.Parameter.Value + " | continuing")
	}

	queryParams := req.QueryStringParameters
	action := queryParams["action"]
	secret := queryParams["secret"]
	test := queryParams["test"]
	body := queryParams["body"]

	if action == "" {
		slog.Error("no action")
		return "", nil
	}
	if action == "manual" && secret != os.Getenv("TOKEN") {
		slog.Error("unauthorized")
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
		slog.Info("manual message " + channel + " " + body)
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

	slog.Info("relative = " + timestampRel)
	slog.Info("full = " + timestampFull)

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rand.Intn(len(imageUrls))
	gifUrl := imageUrls[randomIndex]

	slog.Info("random gif URL = " + gifUrl)

	if now.Year() == threeDaysBefore.Year() && now.Month() == threeDaysBefore.Month() && now.Day() == threeDaysBefore.Day() {
		slog.Info("3 days")
		_, err2 := bot.ChannelMessageSend(channel, "@everyone next session in "+timestampRel+" ("+timestampFull+")")
		check(err2)
	} else if now.Year() == oneDayBefore.Year() && now.Month() == oneDayBefore.Month() && now.Day() == oneDayBefore.Day() {
		slog.Info("1 day")
		_, err2 := bot.ChannelMessageSend(channel, "@everyone session "+timestampRel)
		check(err2)
	} else if now.Year() == secondSaturday.Year() && now.Month() == secondSaturday.Month() && now.Day() == secondSaturday.Day() && now.Hour() == oneHourBefore.Hour() {
		slog.Info("1 hour")
		_, err2 := bot.ChannelMessageSend(channel, "@everyone [session]("+gifUrl+") in "+timestampRel)
		check(err2)
	} else {
		slog.Info("no reminder today")
	}

	return "", nil
}
