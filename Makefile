DOCKER_NETWORK=test-network
DOCKER_VOLUME=test-volume
APP_CONTAINER_NAME=habit-bot
APP_IMAGE_NAME=ghcr.io/alexadastra/habit-bot/${APP_CONTAINER_NAME}
DB_CONTAINER_NAME=habit-bot-mongo
REDIS_CONTAINER_NAME=habit-bot-redis

ifneq (,$(wildcard ./values.env))
    include values.env
    export
endif

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
	@ docker run --rm -d \
		--network ${DOCKER_NETWORK} \
		--env-file values.env \
		--volume ${DOCKER_VOLUME}:/data/db \
		--name ${DB_CONTAINER_NAME} mongo:4.4
	@ docker network inspect -f '{{json .Containers}}' ${DOCKER_NETWORK} | \
		jq '.[] | "Container " + .Name + " runs at ip: " + .IPv4Address'

mongo-stop:
	@ docker stop ${DB_CONTAINER_NAME} || true

mongo-rm: mongo-stop
	@ docker rm ${DB_CONTAINER_NAME} || true

migrations-up:
	@ go run cmd/migrations/main.go up

redis-run:
	@ docker run --rm -d \
		--network ${DOCKER_NETWORK} \
		-p 6379:6379 \
		--name ${REDIS_CONTAINER_NAME} \
		redis:7.2-rc1 redis-server \
		--requirepass ${REDIS_PASSWD}

redis-stop:
	@ docker stop ${REDIS_CONTAINER_NAME}

app-build:
	docker build -t ${APP_IMAGE_NAME} .

app-push:
	docker push ${APP_IMAGE_NAME}

app-pull:
	docker pull ${APP_IMAGE_NAME}

app-run:
	@ docker run --rm -d \
		--network ${DOCKER_NETWORK} \
		--env-file values.env \
		-p 8080:8080 \
		--name ${APP_CONTAINER_NAME} ${APP_IMAGE_NAME}:latest

app-stop:
	@ docker stop ${APP_CONTAINER_NAME} || true

app-rm: app-stop
	@ docker rm ${APP_CONTAINER_NAME} || true

dev-run: app-build network-create volume-create mongo-run migrations-up app-run

dev-stop: app-stop redis-stop mongo-stop
	@ echo "Stopped successfully!"

dev-clean: app-rm mongo-rm volume-rm network-rm
	@ echo "Cleaned up successfully!"