BOT-TOKEN?=
MONGO-DB-DSN?=

build:
	go build -o main .

run:
	go run main.go

test:
	go test -v ./...

docker-build:
	docker build -t habit-bot .

docker-run:
	docker run habit-bot:latest -e "BOT_TOKEN=${BOT-TOKEN}" -e "MONGO_DB_DSN={MONGO-DB-DSN}"

docker-run-env:
	docker run -d habit-bot -p 8081:8081 -e BOT-TOKEN=${BOT-TOKEN} -e MONGO-DB-DSN=${MONGO-DB-DSN}