#!/usr/bin/env bash
docker rm -f gatewayinfo344
docker rm -f usersdb
docker rm -f redissvr

docker pull aaronluannguyen/gatewayinfo344
docker run -it --rm

docker network rm aaronchatnet
docker network create aaronchatnet

export MYSQL_ROOT_PASSWORD=qejblkafuhqieuhgqeu1239
export MYSQL_DATABASE=users
export MYSQL_ADDR=usersdb:3306
export REDISADDR=redissvr:6379
export SESSIONKEY=qlwfnfvdfvkbubiu9859b
export DSN="root:$MYSQL_ROOT_PASSWORD@tcp($MYSQL_ADDR)/$MYSQL_DATABASE"

docker pull aaronluannguyen/usersdb
docker run -d \
--network aaronchatnet \
--name usersdb \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
aaronluannguyen/usersdb


docker run -d \
--name redissvr \
--network aaronchatnet \
redis

docker run -d --name gatewayinfo344 \
--network aaronchatnet \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=/etc/letsencrypt/live/api.aaronluannguyen.me/fullchain.pem \
-e TLSKEY=/etc/letsencrypt/live/api.aaronluannguyen.me/privkey.pem \
-e DSN=$DSN \
-e SESSIONKEY=$SESSIONKEY \
-e REDISADDR=$REDISADDR \
aaronluannguyen/gatewayinfo344