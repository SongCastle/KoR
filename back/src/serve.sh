#!/bin/sh

chmod +x build.sh migrate.sh

./build.sh || exit 1
./migrate.sh up || exit 1

/go/bin/app
