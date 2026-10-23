# Build the application
FROM golang:1.25.5-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP=api

ENV CGO_ENABLED=0
RUN go build -o bin/service ./cmd/${APP}

# Run the tests
FROM build AS run-test
RUN go test -v ./...

# Run the binary
FROM alpine:latest AS run

WORKDIR /app

COPY --from=build /app/bin /app

ENTRYPOINT ["./service"]
