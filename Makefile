tests:
	go test ./internal/... -race -count 10

build:
	docker build .

run: build
	docker-compose up

run-composer-build:
	docker-compose up --build
