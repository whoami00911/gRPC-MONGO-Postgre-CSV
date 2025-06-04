.DEFAULT_GOAL := start

compose-up:
	docker compose up 

rebuild:
	docker compose up --build

down:    
	docker compose down

build:
	docker compose build

restart:    
	docker compose down
	docker compose up

testing:
	go test --short -coverprofile=coverage.out -v ./...    
	go tool cover -func=coverage.out | grep "total"

start: compose-up