#!/usr/bin/env bash
./build.sh
docker push aaronluannguyen/clientinfo344
ssh root@aaronluannguyen.me 'bash-s' < update-client.sh