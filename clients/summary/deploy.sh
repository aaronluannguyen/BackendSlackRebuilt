#!/usr/bin/env bash
echo "Running build script..."
./build.sh
docker push aaronluannguyen/summary-client
ssh root@138.197.220.38 'bash -s' < update-client.sh