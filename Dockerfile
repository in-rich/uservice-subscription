FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o /server cmd/server/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /server /server

ARG PORT=8080
ARG DSN
ARG ENV

ENV PORT=$PORT
ENV DSN=$DSN
ENV ENV=$ENV

ENV HOST="0.0.0.0"

EXPOSE $PORT

# Run
CMD ["/server"]
