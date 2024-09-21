#!/bin/bash

exec_sql() {
  local sql=""
  local db=$1
  local user=$2

  if [ -n "$db" ]; then
    sql+="CREATE DATABASE IF NOT EXISTS \`$db\`;"

    if [ -n "$user" ]; then
      sql+="GRANT ALL ON \`$db\`.* TO '$user'@'%';"
    fi
  fi

  if [ -n "$sql" ]; then
    mysql -h localhost -uroot -p"$MYSQL_ROOT_PASSWORD" <<< "$sql"
  fi
}

exec_sql "$MYSQL_DATABASE" "$MYSQL_USER"

exec_sql "$MYSQL_DATABASE_TEST" "$MYSQL_USER"
