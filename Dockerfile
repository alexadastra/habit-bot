############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
# RUN apk update && apk add --no-cache git
WORKDIR /go/src/github.com/alexadastra/habit-bot
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d ./...
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/habit-bot ./cmd/habit-bot
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/habit-bot /go/bin/habit-bot
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Fetch env
ENV BOT-TOKEN=${BOT-TOKEN}
ENV MONGO-DB-DSN=${MONGO_DB_DSN}
# Run the habit-bot binary.
ENTRYPOINT ["/go/bin/habit-bot"]
