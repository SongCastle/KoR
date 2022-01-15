#!/bin/sh

chmod +x build.sh migrate.sh

./build.sh || exit 1
./migrate.sh up || exit 1

MYSQL_PASSWORD=$(cat $MYSQL_ROOT_PASSWORD_FILE) /go/bin/app
