#!/usr/bin/env bash
echo "Building Linux Executable..."
GOOS=linux go build
echo "Building Docker Container Image..."
docker build -t aaronluannguyen/gatewayinfo344 .
echo "Cleaning Up..."
go clean
docker image prune -f