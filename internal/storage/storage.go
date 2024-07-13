package storage

import (
	"context"
	"event_ticket/internal/model"
)

type User interface {
	CreateUser(ctx context.Context, usr model.CreateUserRequest) (model.User, error)
	GetUser(ctx context.Context, username string) (model.User, error)
}

type Event interface {
	PostEvent(ctx context.Context, eventParam model.Event) (model.Event, error)
	FetchEvents(ctx context.Context) ([]model.Event, error)
	FetchEvent(ctx context.Context, id int32) (model.Event, error)
}

type Payment interface {
	RecordPaymentIntent(ctx context.Context, param model.CreateIntentParam) (model.Payment, error)
	GetPayment(ctx context.Context, intent_id string) (model.Payment, error)
}

type Ticket interface {
	HoldTicket(ctx context.Context, req model.ReserveTicketRequest) (model.Ticket, error)
	GetTicket(id string) (model.Ticket, error)
	UnholdTicket(tktNo, tripID int32) (model.Ticket, error)
}
