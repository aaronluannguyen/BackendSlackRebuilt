#!/usr/bin/env bash
echo "Running build script..."
./build.sh
docker push aaronluannguyen/summarymicroservice
ssh root@api.aaronluannguyen.me 'bash -s' < update-server.sh