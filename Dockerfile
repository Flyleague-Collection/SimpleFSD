FROM golang:1.24.6-alpine AS builder

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o /app/simple-fsd ./cmd/fsd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /fsd

COPY --from=builder /app/simple-fsd .

EXPOSE 6809

ENTRYPOINT ["./simple-fsd"]