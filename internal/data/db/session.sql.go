package db

import (
	"context"
	"time"
)

const storeCheckoutSession = `-- name: StoreCheckoutSession :one
	INSERT INTO sessions(id,ticket_id,payment_url,cancel_url,payment_status,created_at,amount)Values($1,$2,$3,$4,$5,$6,$7)
`

type StoreCheckoutSessionParams struct {
	ID            string
	TicketID      string
	PaymentStatus string
	PaymentURL    string
	CancelURL     string
	Amount        float64
	CreatedAt     time.Time
}

func (q *Queries) StoreCheckoutSession(ctx context.Context, arg StoreCheckoutSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, storeCheckoutSession,
		arg.ID,
		arg.TicketID,
		arg.PaymentStatus,
		arg.PaymentURL,
		arg.CancelURL,
		arg.Amount,
		arg.CreatedAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.TicketID,
		&i.PaymentStatus,
		&i.PaymentURL,
		&i.CancelURl,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTicketStatus = `-- name: GetCheckoutSession :one 
	SELECT status FROM tickets t
	JOIN sessions s ON (t.id=s.ticket_id) 
	WHERE s.id = $1;
`

func (q *Queries) GetTicketStatus(ctx context.Context, sid string) (string, error) {
	row := q.db.QueryRow(ctx, getTicketStatus, sid)

	var i string

	err := row.Scan(
		&i,
	)
	if err != nil {
		return "", err
	}
	return i, nil
}
