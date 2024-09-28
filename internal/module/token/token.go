package token

import (
	"context"
	"database/sql"
	"errors"
	"event_ticket/internal/data/db"
	"fmt"
	"log/slog"
)

type token struct {
	log *slog.Logger
	db.Querier
}

type UnknownServerError struct {
	message string
	cause   error
}

func (e *UnknownServerError) Error() string {
	return e.message
}

type NotFoundError struct {
	resourceName string
	resourceId   string
	context      string
}

func (t *NotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find resource %s with id %s while %s", t.resourceName, t.resourceId, t.context)
}

func (t *token) GenerateToken(ctx context.Context, payload TokenPayload) error {
	ticketInfo, err := t.GetTicket(ctx, payload)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &NotFoundError{
				resourceName: "ticket",
				resourceId:   payload.Id,
				context:      "generating token",
			}
		}
		return &UnknownServerError{
			message: "Unknown server error happened while fetching ticket data",
			cause:   err,
		}
	}

	if ticketInfo.Status == "Onhold" {
		return fmt.Errorf(`Ticket status not updated yet`)
	}

	if ticketInfo.Status != "Reserved" {
		return fmt.Errorf(`Application returned with invalid ticket status`)
	}

	data, err := t.GetTokenData(payload.sessionId); err != nil {
		return fmt.Errorf("Couldn't fetch token data")
	}
	
	return generateToken(data)
}

func generateToken(data any) {
	// getSigningKey
	// know signing algorithms
	// sign the token and return
}

// ticket info
// user info

// generateInfoForToken() {}
// qrcode metadata store, query
// ui graphics metadata store marege query
// UI upload mareg
