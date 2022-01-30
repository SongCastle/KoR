#!/bin/sh

cmd=$1

pass=$(cat $MYSQL_ROOT_PASSWORD_FILE)
db="mysql://$MYSQL_USERNAME:$pass@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=True"

migrate -source file://migrations/ -database $db $cmd
