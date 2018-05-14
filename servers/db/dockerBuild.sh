#!/usr/bin/env bash
export MYSQL_ROOT_PASSWORD="password"
export MYSQL_DATABASE=users

docker rm -f usersdb
docker run -d \
--network aaronchatnet \
--name usersdb \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
aaronluannguyen/usersdb