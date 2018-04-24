#!/usr/bin/env bash

# ensure MYSQL_ROOT_PASSWORD
# and MYSQL_DATABASE are set
if [[ -z $MYSQL_ROOT_PASSWORD ]]
then
    echo "please set MYSQL_ROOT_PASSWORD"
    exit 1
fi
if [[ -z $MYSQL_DATABASE ]]
then
    echo "please set MYSQL_DATABASE"
    exit 1
fi

CONTAINER_NAME=usersdb

# stop and remove any existing container instance
if [ "$(docker ps -aq --filter name=$CONTAINER_NAME)" ]
then
    docker rm -f $CONTAINER_NAME
fi

# run the container
docker build -t aaronluannguyen/usersdb .
docker push aaronluannguyen/usersdb
docker rm -f usersdb
docker pull aaronluannguyen/usersdb
docker network rm aaronchatnet
docker network create aaronchatnet
docker run -d \
--network aaronchatnet \
--name $CONTAINER_NAME \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
aaronluannguyen/usersdb