package initiator

import (
	"context"
	"event_ticket/internal/data/db"
	"log"

	"github.com/jackc/pgx/v4"
)

func InitDB(connStr string) *db.Queries {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db.New(conn)
}
