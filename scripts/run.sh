#!/bin/bash

int_handler()
{
    docker compose down
}
trap int_handler INT

docker compose up -d

go run cmd/server/main.go
