
build:
	go build -o main .

run:
	go run main.go

test:
	go test -v ./...

docker-build:
	docker build -t habit-bot .

docker-run:
	docker run --env-file values.env habit-bot:latest