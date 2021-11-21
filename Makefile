tests:
	go test -race -count 100

run: build
	docker-compose up --build server
