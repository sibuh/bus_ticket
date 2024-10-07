package module

import (
	"context"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"

	"github.com/google/uuid"
)

type Ticket interface {
	ReserveTicket(ctx context.Context, req model.ReserveTicketRequest) (db.Session, error)
}

type User interface {
	CreateUser(ctx context.Context, usr model.CreateUserRequest) (db.User, error)
	LoginUser(ctx context.Context, logReq model.LoginRequest) (string, error)
	RefreshToken(ctx context.Context, username string) (string, error)
}
type Event interface {
	PostEvent(ctx context.Context, postEvent model.Event) (model.Event, error)
	FetchEvents(ctx context.Context) ([]model.Event, error)
	FetchEvent(ctx context.Context, id int32) (model.Event, error)
}

//	type Payment interface {
//		CreatePaymentIntent(ctx context.Context, userID, eventID int32) (string, error)
//		// GetPayment(ctx context.Context, intentID string) (db.Payment, error)
//	}
type Token interface {
	GenerateToken(ctx context.Context, tid, uid uuid.UUID) (string, error)
}
