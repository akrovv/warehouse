.PHONY: up down test lint

up:
	docker-compose up

down:
	docker-compose down

test:
	go test --cover ./... | grep -v 'no test files'

lint:
	golangci-lint -c golangci.yml run ./...
