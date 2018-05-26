#!/usr/bin/env bash
echo "Building Docker Container Image..."
docker build -t aaronluannguyen/messagesmicroservice .
echo "Cleaning Up..."
go clean
docker image prune -f