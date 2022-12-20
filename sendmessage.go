package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/satori/go.uuid"
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

func SendMessage(sess *session.Session, queueUrl string, messageBody string) error {
	sqsClient := sqs.New(sess)
	groupid := "abd" //uuid.NewV4().String()
	input := &sqs.SendMessageInput{
		QueueUrl:       &queueUrl,
		MessageBody:    aws.String(messageBody),
		MessageGroupId: &groupid,
	}
	input.SetMessageDeduplicationId(uuid.NewV4().String())

	_, err := sqsClient.SendMessage(input)

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
	for i := 1; i <= 10; i++ {
		messageBody := "This is a test message" + fmt.Sprintf("%v", i)
		err = SendMessage(sess, *urlRes.QueueUrl, messageBody)
		if err != nil {
			fmt.Printf("Got an error while trying to send message to queue: %v", err)
			return
		}
		log.Println(messageBody)
	}

	fmt.Println("Message sent successfully")
}
