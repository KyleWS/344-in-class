#!/usr/bin/env bash

# you should `source` this script before building the
# Docker container image or trying to run it. Use this
# command to source this into your current shell:
#   source devenv.sh

# docker container image name for our customized MySQL image
# TODO: change the value to use your DockerHub user name
export TASKS_MYSQL_IMAGE=drstearns/mysqltasks

# database name in which our schema will be created
export MYSQL_DATABASE=tasks

# random MySQL root password
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
