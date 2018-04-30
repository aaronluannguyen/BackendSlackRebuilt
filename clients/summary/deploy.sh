#!/usr/bin/env bash
echo "Running build script..."
./build.sh
docker push aaronluannguyen/summary-client
ssh root@aaronnluannguyen.me 'bash -s' < update-client.sh