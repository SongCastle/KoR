#!/bin/bash

if [ -z "$MYSQL_DATABASE" ] || [ -z "$MYSQL_DATABASE_TEST" ]; then
  echo 'Empty env variables (Database)'
  exit 1
fi

mysql=(mysql -uroot -p${MYSQL_ROOT_PASSWORD})
"${mysql[@]}" << SQL
CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE;
CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE_TEST;
SQL

echo 'Create Database'
