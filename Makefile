build:
	@go build -o bin/tireapp cmd/main.go

test:
	@go test -v ./...
	
run: build
	@./bin/tireapp

migration:
	@migrate create -ext sql -dir db/migrate/migrations/ -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run db/migrate/main.go up 

migrate-down:
	@go run db/migrate/main.go down 
