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
ENV MONGODB_USER=value2
ENV MONGODB_PASSWORD=value3
ENV MONGODB_DATABASE=value4
ENV MONGODB_HOST=value5

# Fetch dependencies using go get.
RUN go get -d ./... && \
# Build the binary.
CGO_ENABLED=0 go build -ldflags="-w -s" -o main ./cmd/habit-bot
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