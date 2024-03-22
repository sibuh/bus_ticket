package ticket

import (
	"context"
	"event_ticket/internal/data/db"
	"event_ticket/internal/model"
	"event_ticket/internal/storage"
	"log"

	"github.com/jackc/pgx/v4"
	"golang.org/x/exp/slog"
)

type ticket struct {
	log     slog.Logger
	queries *db.Queries
}

func Init(logger slog.Logger, constr string) storage.Ticket {
	conn, err := pgx.Connect(context.Background(), constr)
	if err != nil {
		logger.Error("failed to creat database connection object", err)
		log.Fatal(err)

	}
	queries := db.New(conn)
	return &ticket{
		log:     logger,
		queries: queries,
	}
}
func (t *ticket) RegisterUserToDb(user model.User, sessionid, nonce string) error {
	_, err := t.queries.RegisterPayedUser(context.Background(), db.RegisterPayedUserParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Nonce:     nonce,
		SessionID: sessionid,
	})
	if err != nil {
		t.log.Error("failed to register user", err)
		return err
	}
	return nil

}

func (t *ticket) UpdatePaymentStatus(status, sid string) (db.User, error) {
	user, err := t.queries.UpdatePaymentStatus(context.Background(), db.UpdatePaymentStatusParams{
		PaymentStatus: status,
		SessionID:     sid,
	})
	if err != nil {
		t.log.Error("failed to update the payment status", err)
		return db.User{}, err
	}
	return user, nil
}
func (t *ticket) GetUser(nonce string) (db.User, error) {
	user, err := t.queries.GetUser(context.Background(), nonce)
	if err != nil {
		t.log.Error("failed to get user with the given nonce", err)
		return db.User{}, err
	}
	return user, nil

}
