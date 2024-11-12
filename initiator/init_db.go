package initiator

import (
	"bus_ticket/internal/data/db"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDB(connStr string) *db.Queries {

	conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db.New(conn)
}
