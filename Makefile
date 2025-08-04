.PHONY: build run test clean docker-build docker-run

build:
	go build -o bin/secret-manager .

run:
	go run .

test:
	go test ./tests/... -v

clean:
	rm -rf bin/

docker-build:
	docker build -t secret-manager .

docker-run:
	docker-compose up --build

docker-down:
	docker-compose down

lint:
	golangci-lint run

deps:
	go mod tidy
	go mod download

migrate:
	go run . migrate

dev:
	air
