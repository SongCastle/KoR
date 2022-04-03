## 技術スタック
golang ([Gin](https://github.com/gin-gonic/gin), [gorm](https://github.com/go-gorm/gorm), [golang-migrate](https://github.com/golang-migrate/migrate)), MySQL

### API 一覧

```
GET    /ping           # 疎通確認用
GET    /v1/users       # ユーザ一覧取得
GET    /v1/users/:id   # ユーザ取得
PUT    /v1/users/:id   # ユーザ更新
POST   /v1/users       # ユーザ作成
DELETE /v1/users/:id   # ユーザ削除
POST   /v1/users/token # ユーザ認証
DELETE /v1/users/token # ユーザ認証無効化
POST   /v1/users/admin/token # ユーザ認証 (管理者)
```

#### 操作例

##### 疎通確認
```
$ curl 0.0.0.0:3000/ping
> pong
```

##### ユーザ操作
```
# 作成
$ curl -X POST 0.0.0.0:3000/v1/users -H 'Content-Type: application/json' -d '{"login": "user", "password": "user1234", "email": "user@example.com"}'
> {"id":1,"login":"user","email":"user@example.com","created_at":"2022-01-15T10:50:20.3272622Z","updated_at":"2022-01-15T10:50:20.3272622Z"}

# 取得
$ curl 0.0.0.0:3000/v1/users/1
> {"id":1,"login":"user","email":"user@example.com","created_at":"2022-01-15T10:50:20.3272622Z","updated_at":"2022-01-15T10:50:20.3272622Z"}

# 削除
$ curl -X DELETE -H 'Authorization: Bearer <token>' 0.0.0.0:3000/v1/users/1

# 認証
$ curl -X PUT -H 0.0.0.0:3000/v1/users/token -H 'Content-Type: application/json' -d '{"login": "user", "password": "user1234"}'
> <token>
```
