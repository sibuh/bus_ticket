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
