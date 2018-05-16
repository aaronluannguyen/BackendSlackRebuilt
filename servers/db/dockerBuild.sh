#!/usr/bin/env bash
export MYSQL_ROOT_PASSWORD="password"
export MYSQL_DATABASE=users

docker build -t aaronluannguyen/usersdb .
docker rm -f aaronluannguyen/usersdb
docker push aaronluannguyen/userdb
docker pull aaronluannguyen/usersdb
docker run -d \
--network aaronchatnet \
--name usersdb \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
aaronluannguyen/usersdb