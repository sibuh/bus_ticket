migrate-up:
	-migrate -path internal/data/migrations -database "cockroachdb://root@localhost:26257/tickets?sslmode=disable" -verbose up
migrate-down:
	-migrate -path internal/data/migrations -database "cockroachdb://root@localhost:26257/tickets?sslmode=disable" -verbose down
sqlc:
	- sqlc generate -f ./config/sqlc.yaml
run:
	-go run cmd/main.go