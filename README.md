# KoR
WORK IN PROGRESS ...

# Getting Started

1. 本レポジトリを Clone します
```
$ git clone git@github.com:SongCastle/KoR.git
$ cd KoR
```

2. 環境変数関連のファイルをコピーします
```
$ cp back/.env.sample back/.env
$ cp db/.env.sample db/.env
$ cp db/password.txt.sample db/password.txt
$ cp front/.env.sample front/.env
```

3. コンテナを起動します
```
$ docker-compose up -d
```

4. back コンテナに入り、WEB サーバを起動します
```
# コンテナ外
$ docker-compose exec back ash

# コンテナ内
# chmod +x serve.sh && ./serve.sh
```
