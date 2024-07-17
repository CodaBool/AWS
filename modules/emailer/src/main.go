package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	_ "github.com/joho/godotenv/autoload"
)

type Input struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func main() {
	buildLogger()
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		handle(context.TODO(), Input{
			Subject: "sub",
			Message: "msg",
		})
	} else {
		lambda.Start(handle)
	}
}

func handle(ctx context.Context, input Input) (string, error) {
	log.Print("input ", input.Subject, input.Message)
	log.Print("arn ", os.Getenv("TOPIC_ARN"))
	cfg, err := config.LoadDefaultConfig(ctx)
	check(err)
	client := sns.NewFromConfig(cfg)
	res, err := client.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(input.Message),
		Subject:  aws.String(input.Subject),
		TopicArn: aws.String(os.Getenv("TOPIC_ARN")),
	})
	check(err)
	log.Print(res)
	return "email sent", nil
}
