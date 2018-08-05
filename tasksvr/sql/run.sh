#!/usr/bin/env bash
set -e

CONTAINER_NAME=tasksdb

./build.sh

if [ "$(docker ps -aq --filter name=$CONTAINER_NAME)" ]; then
    docker rm -f $CONTAINER_NAME
fi

docker run -d \
--name $CONTAINER_NAME \
-p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
$TASKS_MYSQL_IMAGE --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
