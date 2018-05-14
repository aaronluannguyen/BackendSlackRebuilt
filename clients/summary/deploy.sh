#!/usr/bin/env bash
echo "Running build script..."
./build.sh
docker push aaronluannguyen/summary-client
ssh root@206.189.64.90 'bash -s' < update-client.sh