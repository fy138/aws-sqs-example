package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func DeleteMessage(sess *session.Session, queueUrl string, messageHandle *string) error {
	sqsClient := sqs.New(sess)

	_, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: messageHandle,
	})

	return err
}
func GetQueueURL(sess *session.Session, queue string) (*sqs.GetQueueUrlOutput, error) {
	sqsClient := sqs.New(sess)

	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Println("env1", os.Getenv("AWS_ACCESS_KEY_ID"))
	log.Println("env2", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	log.Println("env3", os.Getenv("AWS_REGION"))
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})

	if err != nil {
		fmt.Printf("Failed to initialize new session: %v", err)
		return nil
	}
	/*
		queueName := "myqueueone.fifo"

		urlRes, err := GetQueueURL(sess, queueName)
		if err != nil {
			fmt.Printf("GetQueueURL error: %v", err)
			return nil
		}
		urlRes.QueueUrl
	*/
	qurl := `https://sqs.us-east-1.amazonaws.com/843344726074/myqueueone.fifo`
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.ReceiptHandle, message.EventSource, message.Body)
		err = DeleteMessage(sess, qurl, &message.ReceiptHandle)
		if err != nil {
			fmt.Printf("Got an error while trying to delete message: %v", err)
			return nil
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
