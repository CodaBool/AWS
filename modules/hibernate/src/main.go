package main

import (
	"context"
	"fmt"

	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Input struct {
	Start bool `json:"start"`
}

type Instance struct {
	Running   bool
	Hibernate bool
	Name      string
	Id        string
}

func main() {
	buildLogger()
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		handle(context.TODO(), Input{Start: true})
	} else {
		lambda.Start(handle)
	}
}

func handle(ctx context.Context, input Input) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	check(err)
	clientEC2 := ec2.NewFromConfig(cfg)
	res, err := clientEC2.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	check(err)
	var instances []Instance

	for _, reservation := range res.Reservations {
		for _, instance := range reservation.Instances {
			if instance.State.Name == "terminated" {
				continue
			}
			log.Print("state=", instance.State.Name)
			var i Instance
			if instance.State.Name == "stopped" {
				i.Running = false
			} else if instance.State.Name == "running" {
				i.Running = true
			}
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					log.Print("name=", *tag.Value)
					i.Name = *tag.Value
				}
				if *tag.Key == "hibernate" {
					log.Print("hibernate=", *tag.Value)
					if *tag.Value == "true" {
						i.Hibernate = true
					} else {
						i.Hibernate = false
					}
				}
			}
			i.Id = *instance.InstanceId
			instances = append(instances, i)
		}
	}

	var ids []string
	for _, i := range instances {
		ids = append(ids, i.Id)
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
		_, err = clientEC2.StartInstances(ctx, &ec2.StartInstancesInput{
			InstanceIds: ids,
		})
		check(err)
	} else {
		subject = "Stopped " + plural_inst
		_, err = clientEC2.StopInstances(ctx, &ec2.StopInstancesInput{
			InstanceIds: ids,
		})
		check(err)
	}
	log.Print(subject)

	// message := strings.Join(ids[:], ", ")
	// clientSNS := sns.NewFromConfig(cfg)
	// _, err = clientSNS.Publish(ctx, &sns.PublishInput{
	// 	Message:  aws.String(message),
	// 	Subject:  aws.String(subject),
	// 	TopicArn: aws.String("arn:aws:sns:us-east-1:919759177803:notify"),
	// })
	// check(err)

	return subject, nil
}
