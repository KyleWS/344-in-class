#!/usr/bin/env bash
docker run -it \
--network host \
--rm \
mysql sh -c "exec mysql -h127.0.0.1 -uroot -p$MYSQL_ROOT_PASSWORD"
