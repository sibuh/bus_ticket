package token

import (
	"bus_ticket/internal/constant"
	"bus_ticket/internal/data/db"
	"bus_ticket/internal/module"
	tkn "bus_ticket/internal/utils/token"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type token struct {
	log *slog.Logger
	db.Querier
	paseto tkn.TokenMaker
}

var ErrInvalidPayload = errors.New("invalid payload")

func Init(log *slog.Logger, q db.Querier, maker tkn.TokenMaker) module.Token {
	return &token{
		log:     log,
		Querier: q,
		paseto:  maker,
	}
}

func (t *token) GenerateToken(ctx context.Context, tid, uid uuid.UUID) (string, error) {
	ticketInfo, err := t.Querier.GetTicket(ctx, tid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", &NotFoundError{
				resourceName: "ticket",
				resourceId:   tid.String(),
				context:      "generating token",
			}
		}
		t.log.Error("db returned error while fetching ticket", tid, ctx, uid)
		return "", &UnknownServerError{
			message: "Unknown server error happened while fetching ticket data",
			cause:   err,
		}
	}

	if ticketInfo.Status == "Onhold" {
		// t.log.Info("generate token request before ticket status is updated")
		return "", &ErrStatusNotUpdated{ID: tid.String(), Status: string(constant.Onhold), Retry: true}
	}

	if ticketInfo.Status != "Reserved" {
		t.log.Error("Got invalid ticket status", ErrInvalidTicketStatus{
			ID:      tid.String(),
			Status:  ticketInfo.Status,
			Message: fmt.Sprintf("got invalid ticket status %s expected ticket status %s", ticketInfo.Status, constant.Reserved),
		})
		return "", &UnknownServerError{
			message: "Unknown server error happened",
			cause:   err,
		}
	}
	// duration := time.
	return t.generateToken(tid, uid, 24*time.Hour)
}

func (t *token) generateToken(tid, uid uuid.UUID, duration time.Duration) (string, error) {
	ticketPayload := NewTicketTokenPayload(uid, tid, duration)
	return t.paseto.CreateToken(ticketPayload)
}
func (t *token) VerifyToken(ctx context.Context, tkn string) (TicketInfo, error) {
	payload := Payload{}
	pl, err := t.paseto.VerifyToken(tkn, &payload)
	if err != nil {
		return TicketInfo{}, err
	}
	p, ok := pl.(*Payload)
	if !ok {
		return TicketInfo{}, ErrInvalidPayload
	}
	tkt, err := t.Querier.GetTicket(ctx, p.TicketID)
	if err != nil {
		return TicketInfo{}, err
	}

	u, err := t.Querier.GetUser(ctx, p.UserID.String())
	if err != nil {
		return TicketInfo{}, err
	}
	return TicketInfo{u, tkt}, nil
}
