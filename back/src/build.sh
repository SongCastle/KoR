#!/bin/sh

cd cmd/app
go build -o /go/bin/app

if [ ! $? = 0 ]; then
  echo 'build filed'
  exit 1
fi
