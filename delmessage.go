package main

import (
	"fmt"
	"log"

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

func GetMessages(sess *session.Session, queueUrl string, maxMessages int64) (*sqs.ReceiveMessageOutput, error) {
	sqsClient := sqs.New(sess)

	msgResult, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(maxMessages),
	})

	if err != nil {
		return nil, err
	}

	return msgResult, nil
}
func DeleteMessage(sess *session.Session, queueUrl string, messageHandle *string) error {
	sqsClient := sqs.New(sess)

	_, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: messageHandle,
	})

	return err
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

	queueName := "myqueueone.fifo"

	urlRes, err := GetQueueURL(sess, queueName)
	if err != nil {
		fmt.Printf("Got an error while trying to create queue: %v", err)
		return
	}

	var maxMessages int64 = 3
	msgRes, err := GetMessages(sess, *urlRes.QueueUrl, maxMessages)
	if err != nil {
		fmt.Printf("Got an error while trying to retrieve message: %v", err)
		return
	}
	for _, li := range msgRes.Messages {
		//log.Printf("%#v", msgRes)

		receiptHandle := li.ReceiptHandle
		err = DeleteMessage(sess, *urlRes.QueueUrl, receiptHandle)
		if err != nil {
			fmt.Printf("Got an error while trying to delete message: %v", err)
			return
		}
		log.Println("Message Body: " + *li.Body)
		log.Println("Deleted message with handle: " + *receiptHandle)
	}

}
