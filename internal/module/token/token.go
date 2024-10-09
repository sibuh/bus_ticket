package token

import (
	"context"
	"database/sql"
	"errors"
	"event_ticket/internal/constant"
	"event_ticket/internal/data/db"
	"event_ticket/internal/module"
	tkn "event_ticket/internal/utils/token"
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

func Init(log *slog.Logger, q db.Querier, tkn tkn.TokenMaker) module.Token {
	return &token{
		log:     log,
		Querier: q,
		paseto:  tkn,
	}
}

type UnknownServerError struct {
	message string
	cause   error
}

func (e *UnknownServerError) Error() string {
	return e.message + "cause: " + e.cause.Error()
}

type NotFoundError struct {
	resourceName string
	resourceId   string
	context      string
}
type ErrStatusNotUpdated struct {
	ID     string
	Status string
	Retry  bool
}
type ErrInvalidTicketStatus struct {
	ID      string
	Status  string
	Message string
}

func (e *ErrInvalidTicketStatus) Error() string {
	return e.Message
}

func (e *ErrStatusNotUpdated) Error() string {
	return fmt.Sprintf("status of ticket with id %s is not yet updated.Expected status %s but got %s", e.ID, constant.Reserved, e.Status)
}

func (t *NotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find resource %s with id %s while %s", t.resourceName, t.resourceId, t.context)
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
		return "", &UnknownServerError{
			message: "Unknown server error happened while fetching ticket data",
			cause:   err,
		}
	}

	if ticketInfo.Status == "Onhold" {

		return "", &ErrStatusNotUpdated{ID: tid.String(), Status: string(constant.Onhold), Retry: true}
	}

	if ticketInfo.Status != "Reserved" {
		return "", &ErrInvalidTicketStatus{
			ID:      tid.String(),
			Status:  ticketInfo.Status,
			Message: fmt.Sprintf("got invalid ticket status %s expected ticket status %s", ticketInfo.Status, constant.Reserved),
		}
	}

	return t.generateToken(tid.String(), uid.String(), 24*time.Hour)
}

func (t *token) generateToken(tid, uid string, duration time.Duration) (string, error) {
	userAndTicketID := fmt.Sprintf("%s %s", uid, tid)
	return t.paseto.CreateToken(userAndTicketID, duration)
}
