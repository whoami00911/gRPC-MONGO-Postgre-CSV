.DEFAULT_GOAL := start

compose-up: testing
	docker compose up 

down:    
	docker compose down

rebuild:
	docker compose up --build

build:
	docker compose build

restart:    
	docker compose down
	docker compose up

testing:
	go test --short -coverprofile=coverage.out -v ./...    
	go tool cover -func=coverage.out | grep "total"

start: testing compose-up