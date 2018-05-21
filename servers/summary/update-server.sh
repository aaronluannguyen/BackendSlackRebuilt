#!/usr/bin/env bash
export ADDR=":80"
docker pull aaronluannguyen/summarymicroservice
docker rm -f summarymicroservice
docker run -d \
--network aaronchatnet \
--name summarymicroservice \
-e ADDR=$ADDR \
aaronluannguyen/summarymicroservice