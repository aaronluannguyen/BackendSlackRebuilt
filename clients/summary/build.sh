#!/usr/bin/env bash
echo "Building Docker Container Image..."
docker build -t aaronluannguyen/summary-client .
docker image prune -f