package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSReceiveMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context,
		params *sqs.DeleteMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

func GetQueueURL(c context.Context, api SQSReceiveMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func GetMessages(c context.Context, api SQSReceiveMessageAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}

func RemoveMessage(c context.Context, api SQSReceiveMessageAPI, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return api.DeleteMessage(c, input)
}

func main() {
	queue := flag.String("q", "", "The name of the queue")

	flag.Parse()

	if *queue == "" {
		fmt.Println("You must supply the name of a queue (-q QUEUE)")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sqs.NewFromConfig(cfg)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: queue,
	}

	// Get URL of queue
	urlResult, err := GetQueueURL(context.TODO(), client, gQInput)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := urlResult.QueueUrl

	fmt.Println(*queueURL)

	gMInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: 1,
		//VisibilityTimeout:   int32(1),
	}
	for {
		fmt.Println("No messages found 1")
		time.Sleep(time.Second * 3)
		msgResult, err := GetMessages(context.TODO(), client, gMInput)
		if err != nil {
			fmt.Println("Got an error receiving messages:")
			fmt.Println(err)
			continue
		}
		for _, li := range msgResult.Messages {
			fmt.Println("Message body: " + *li.Body)
			dMInput := &sqs.DeleteMessageInput{
				QueueUrl:      queueURL,
				ReceiptHandle: li.ReceiptHandle,
			}

			_, err = RemoveMessage(context.TODO(), client, dMInput)
			if err != nil {
				fmt.Println("Got an error deleting the message:")
				fmt.Println(err)
				return
			} else {
				fmt.Println("deleting the message success")
			}

		}
	}
}
