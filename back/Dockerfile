FROM golang:1.17.6-alpine

WORKDIR /src/go/github.com/SongCastle/KoR
RUN apk update && apk --no-cache add curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv ./migrate /usr/bin/migrate
