#!/usr/bin/env bash
echo "Building Docker Container Image..."
docker build -t aaronluannguyen/clientinfo344
docker image prune -f