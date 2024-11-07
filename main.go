package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	// tea "github.com/charmbracelet/bubbletea"
	"gomail/app"
	"gomail/service"

	"google.golang.org/api/gmail/v1"
)

func main() {
	app.Run()
}

func startService() {
	ctx := context.Background()
	svc, err := service.NewService(ctx)
	user := "me"
	msgResponse, err := svc.Regular.Users.Messages.List(user).Q(service.NewerThan(1, service.Day).Query).Fields("messages(id,payload/headers)").Do()
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
	_ = msgs

	if err != nil {
		log.Fatalf("Could not marshal message: %v", err)
	}
	// jsonStr := fmt.Sprintf("%s", json)
	path := "./email.json"
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	msg := msgs[0]
	jsonMsg, err := json.Marshal(msg)
	writer.Write(jsonMsg)
	// encoder := json.NewEncoder(file)
	// encoder.Encode(jsonMsg)
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
