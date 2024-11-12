package service

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/api/gmail/v1"
)

type Service struct {
	Regular *gmail.Service
	client  http.Client
	Batch   *BatchEmailService
}

func NewService(ctx context.Context) (*Service, error) {
	gsvc, client, err := GetService(ctx)
	if err != nil {
		return nil, err
	}
	svc := &Service{
		Regular: gsvc,
		client:  *client,
	}
	batch := NewBatchEmailService(svc)
	svc.Batch = batch
	return svc, nil
}

func (svc *Service) RequestRecentMessages(ctx context.Context, user string) []*gmail.Message {
	// retrieve msg ids
	query := NewerThan(1, Day)
	msgIds := svc.requestMsgIds(user, query)
	return svc.requestPayload(ctx, msgIds)
}

func (svc *Service) requestMsgIds(user string, query GmailQuery) []string {
	msgResponse, err := svc.Regular.Users.Messages.List(user).Q(query.Query).Fields("messages(id,payload/headers)").Do()
	// TODO: Error handling
	if err != nil {
		log.Fatalf("Could not retrieve messages: %v", err)
	}
	msgsMeta := msgResponse.Messages
	msgIds := make([]string, 0, len(msgsMeta))
	for _, msg := range msgsMeta {
		msgIds = append(msgIds, msg.Id)
	}
	return msgIds
}

func (svc *Service) requestPayload(ctx context.Context, msgIds []string) []*gmail.Message {
	// retrieve payload
	msgCall := svc.Batch.Get("me", msgIds).Context(ctx).Format("full")
	msgs, err := msgCall.Do()
	// TODO: Error handling
	if err != nil {
		log.Fatalf("Error during batch call: %v", err)
	}
	return msgs
}
