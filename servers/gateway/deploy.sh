#!/usr/bin/env bash
echo "Running build script..."
./build.sh
docker push aaronluannguyen/summary-server
ssh root@api.aaronnluannguyen.me 'bash -s' < update-server.sh