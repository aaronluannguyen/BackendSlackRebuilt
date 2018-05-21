#!/usr/bin/env bash
echo "Running build script..."
./build.sh
docker push aaronluannguyen/messagesmicroservice
ssh root@api.aaronluannguyen.me 'bash -s' < update-server.sh