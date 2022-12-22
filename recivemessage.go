package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

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

func GetMessages(sess *session.Session, queueUrl string, maxMessages int) (*sqs.ReceiveMessageOutput, error) {
	sqsClient := sqs.New(sess)

	msgResult, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(1),
	})

	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func main() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})

	if err != nil {
		fmt.Printf("Failed to initialize new session: %v", err)
		return
	}

	queueName := "myqueuestandard"

	urlRes, err := GetQueueURL(sess, queueName)
	if err != nil {
		fmt.Printf("Got an error while trying to create queue: %v", err)
		return
	}
	log.Println(*urlRes.QueueUrl)
	for {
		time.Sleep(time.Second * 3)
		maxMessages := 1
		msgRes, err := GetMessages(sess, *urlRes.QueueUrl, maxMessages)
		if err != nil {
			fmt.Printf("Got an error while trying to retrieve message: %v", err)
			time.Sleep(time.Second * 10)
			continue
		}
		if len(msgRes.Messages) == 0 {
			continue
		}
		log.Printf("%#v", *msgRes.Messages[0].Body)

	}
	//fmt.Println("Message Body: " + *msgRes.Messages[0].Body)
	//fmt.Println("Message Handle: " + *msgRes.Messages[0].ReceiptHandle)
}
