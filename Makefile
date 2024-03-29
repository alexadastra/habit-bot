DOCKER_NETWORK=test-network
DOCKER_VOLUME=test-volume
DOCKER_APP_NAME=habit-bot
DOCKER_APP_IMAGE=ghcr.io/alexadastra/habit-bot/${DOCKER_APP_NAME}
DOCKER_DB_IMAGE=habit-bot-mongo

build:
	go build -o main .

run:
	go run cmd/habit-bot/main.go

test:
	go test -v ./...

network-create:
	@ docker network create ${DOCKER_NETWORK}

network-rm:
	@ docker network rm ${DOCKER_NETWORK} || true

volume-create:
	@ docker volume create ${DOCKER_VOLUME}

volume-rm:
	@ docker volume rm ${DOCKER_VOLUME} || true

mongo-run:
	@ docker run -d \
		--network ${DOCKER_NETWORK} \
		--env-file values.env \
		--volume ${DOCKER_VOLUME}:/data/db \
		--name ${DOCKER_DB_IMAGE} mongo:4.4
	@ docker network inspect -f '{{json .Containers}}' ${DOCKER_NETWORK} | \
		jq '.[] | "Container " + .Name + " runs at ip: " + .IPv4Address'

mongo-stop:
	@ docker stop ${DOCKER_DB_IMAGE} || true

mongo-rm: mongo-stop
	@ docker rm ${DOCKER_DB_IMAGE} || true

migrations-up:
	@ go run cmd/migrations/main.go up

app-build:
	docker build -t ${DOCKER_APP_IMAGE} .

app-push:
	docker push ${DOCKER_APP_IMAGE}

app-pull:
	docker pull ${DOCKER_APP_IMAGE}

app-run:
	@ docker run -d \
		--network ${DOCKER_NETWORK} \
		--env-file values.env \
		-p 8080:8080 \
		--name ${DOCKER_APP_NAME} ${DOCKER_APP_IMAGE}:latest

app-stop:
	@ docker stop ${DOCKER_APP_NAME} || true

app-rm: app-stop
	@ docker rm ${DOCKER_APP_NAME} || true

dev-run: app-build network-create volume-create mongo-run migrations-up app-run

dev-stop: app-stop mongo-stop
	@ echo "Stopped successfully!"

dev-clean: app-rm mongo-rm volume-rm network-rm
	@ echo "Cleaned up successfully!"