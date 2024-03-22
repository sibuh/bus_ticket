migrate-up:
	-migrate -path data/migrations -database "postgres://hamereeth_yared:Gy~1KSFfuVJ4@hamereth:5432/hamereth_events?sslmode=disable" -verbose up
migrate-down:
	-migrate -path data/migrations -database "postgres://hamereeth_yared:Gy~1KSFfuVJ4@hamereth:5432/hamereth_events" -verbose down
sqlc:
	- sqlc generate -f ./config/sqlc.yaml
run:
	-go run cmd/main.go