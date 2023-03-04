package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Input struct {
	Start bool `json:"start"`
}

func check(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func buildLogger(isProduction bool) {
	if isProduction {
		// has no timestamp and outputs json
		// by default all log levels are printed
		// changing fieldname to message makes filtering easier in AWS
		log.Logger = zerolog.New(os.Stderr).With().Logger()
		// zerolog.SetGlobalLevel(zerolog.InfoLevel)
		// https://github.com/rs/zerolog#error-logging
		zerolog.ErrorFieldName = "message"
		return
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:          os.Stderr,
		PartsExclude: []string{zerolog.TimestampFieldName}, // comment to add time
		FormatCaller: func(i interface{}) string { return "" },
	}).Level(zerolog.DebugLevel).With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {
	prod := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if prod == "" {
		buildLogger(false)
		handle(nil, Input{Start: false})
	} else {
		buildLogger(true)
		lambda.Start(handle)
	}
}

func handle(ctx context.Context, input Input) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	check(err)
	clientEC2 := ec2.NewFromConfig(cfg)
	res, err := clientEC2.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	check(err)
	ids := make([]string, 0)
	// amiIDs := make([]string, 0)
	for _, instance := range res.Reservations[0].Instances {
		// amiIDs = append(amiIDs, *instance.ImageId)
		ids = append(ids, *instance.InstanceId)
	}

	plural_inst := ""
	if len(ids) == 1 {
		plural_inst = "1 instance"
	} else {
		plural_inst = fmt.Sprintf("%v", len(ids)) + " instances"
	}

	subject := ""
	if input.Start {
		subject = "Started " + plural_inst
		_, err = clientEC2.StartInstances(context.TODO(), &ec2.StartInstancesInput{
			InstanceIds: ids,
		})
		check(err)
	} else {
		subject = "Stopped " + plural_inst
		_, err = clientEC2.StopInstances(context.TODO(), &ec2.StopInstancesInput{
			InstanceIds: ids,
		})
		check(err)
	}
	log.Print(subject)

	// message := strings.Join(ids[:], ", ")
	clientSNS := sns.NewFromConfig(cfg)
	_, err = clientSNS.Publish(context.TODO(), &sns.PublishInput{
		Message:  aws.String(subject),
		Subject:  aws.String(subject),
		TopicArn: aws.String("arn:aws:sns:us-east-1:919759177803:notify"),
	})
	check(err)

	return subject, nil
}
