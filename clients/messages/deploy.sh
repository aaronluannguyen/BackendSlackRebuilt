#!/usr/bin/env bash
./build.sh
docker push aaronluannguyen/client344
ssh root@138.68.44.231 'bash -s' < update-client.sh