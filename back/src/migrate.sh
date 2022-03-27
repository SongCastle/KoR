#!/bin/sh

cmd=$1
v=$2

pass=$(cat $MYSQL_ROOT_PASSWORD_FILE)
db="mysql://$MYSQL_USERNAME:$pass@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=True"

if [ -z "$v" ]; then
  migrate -source file://volume/db/migrations/ -database $db $cmd
else
  migrate -source file://volume/db/migrations/ -database $db $cmd $v
fi
