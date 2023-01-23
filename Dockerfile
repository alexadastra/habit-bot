############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS build-stage
# Install git. Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

COPY . /app
WORKDIR /app

# Fetch env
ENV BOT_TOKEN=value1
ENV MONGO_DB_DSN=value2

# Fetch dependencies using go get.
RUN go get -d ./...
# Build the binary.
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o main ./cmd/habit-bot
############################
# STEP 2 build a small image
############################
FROM scratch

# Copy our static executable
COPY --from=build-stage /app/main /
# Pass sertificates
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run the habit-bot binary.
CMD ["./main"]