package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	// "gomail/app"
	"gomail/service"

	"google.golang.org/api/gmail/v1"
)

func main() {
	// app.Run()
	testService()
}

func testService() {
	ctx := context.Background()
	// TODO: fails if token does not work anymore (probably expiry) handle accordingly
	svc, err := service.NewService(ctx)
	if err != nil {
		log.Fatalf("Could not start service %v\n", err)
	}
	user := "me"
	msgs := requestRecentMessages(ctx, user, svc)
	_ = msgs

}

func requestRecentMessages(ctx context.Context, user string, svc *service.Service) []*gmail.Message {
	// retrieve msg ids
	query := service.NewerThan(1, service.Day).Query
	// TODO: implement service wrapper
	msgResponse, err := svc.Regular.Users.Messages.List(user).Q(query).Fields("messages(id,payload/headers)").Do()
	if err != nil {
		log.Fatalf("Could not retrieve messages: %v", err)
	}
	msgsMeta := msgResponse.Messages
	msgIds := make([]string, 0, len(msgsMeta))
	for _, msg := range msgsMeta {
		msgIds = append(msgIds, msg.Id)
	}
	fmt.Printf("Requesting the following message ids: %v\n", msgIds)

	// retrieve payload
	msgCall := svc.Batch.Get("me", msgIds).Context(ctx).Format("full")
	msgs, err := msgCall.Do()
	if err != nil {
		log.Fatalf("Error during batch call: %v", err)
	}
	return msgs
}

func dumpEmail(path string, msg *gmail.Message) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	jsonMsg, err := json.Marshal(msg)
	writer.Write(jsonMsg)
	if err != nil {
		log.Fatalf("Could not marshal message: %v", err)
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
