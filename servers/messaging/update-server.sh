#!/usr/bin/env bash
export MYSQL_ROOT_PASSWORD=qejblkafuhqieuhgqeu1239
export MYSQL_DATABASE=users
export MYSQL_ADDR=usersdb
export ADDR=messagesmicroservice:80
export MQADDR=api.aaronluannguyen.me:5672
export MQNAME=aaronluannguyenMQ

docker pull aaronluannguyen/messagesmicroservice
docker rm -f messagesmicroservice
docker run -d \
--network aaronchatnet \
--name messagesmicroservice \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
-e MYSQL_ADDR=$MYSQL_ADDR \
-e ADDR=$ADDR \
-e MQADDR=$MQADDR \
-e MQNAME=$MQNAME \
aaronluannguyen/messagesmicroservice