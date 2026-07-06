run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

docker-build:
	docker compose build

docker-up:
	docker compose up

docker-down:
	docker compose down

generate:
	./scripts/generate.sh