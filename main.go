package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/gmail/v1"
)

func (q *gmailQuery) append(query string) {
	if q.query == "" {
		q.query = query
	} else {

	}
}

func gmailTime(t time.Time) string {
	return t.Format("04/16/2004")
}

func main() {
	ctx := context.Background()
	svc, err := NewService(ctx)
	user := "me"
	msgResponse, err := svc.Regular.Users.Messages.List(user).Q(NewerThan(1, day).query).Fields("messages(id,payload/headers)").Do()
	if err != nil {
		log.Fatalf("Could not retrieve messages: %v", err)
	}
	msgsMeta := msgResponse.Messages
	msgIds := make([]string, 0, len(msgsMeta))
	for _, msg := range msgsMeta {
		msgIds = append(msgIds, msg.Id)
	}
	fmt.Printf("Requesting the following message ids: %v\n", msgIds)

	msgCall := svc.Batch.Get("me", msgIds).Context(ctx).Format("full")
	msgs, err := msgCall.Do()
	if err != nil {
		log.Fatalf("Error during batch call: %v", err)
	}

	for _, msg := range msgs {
		printMessage(msg)
	}

}

func printMessage(msg *gmail.Message) {
	fmt.Printf("Message id: %s", msg.Id)

	var subject, sender string
	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "Subject":
			subject = header.Value
		case "From":
			sender = header.Value
		}
	}
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("From: %s\n", sender)
}
