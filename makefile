migration-up:
	migrate -path migrations/postgres/ -database "postgresql://zafar:2605@localhost:5432/clinic_exam?sslmode=disable" -verbose up

gen-swag:
	swag init -g api/api.go -o api/docs

run:
	go run cmd/main.go
