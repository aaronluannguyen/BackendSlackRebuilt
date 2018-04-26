#!/usr/bin/env bash
echo "Building Docker Container Image..."
docker build -t aaronluannguyen/client344 .
docker image prune -f