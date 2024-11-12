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
	msgs := svc.RequestRecentMessages(ctx, user)
	_ = msgs
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
