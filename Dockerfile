FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o /server cmd/server/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /server /server

ENV PORT $PORT

ENV DSN $DSN

ENV ENV $ENV

EXPOSE $PORT

ENV HOST "0.0.0.0"

# Run
CMD ["/server"]
