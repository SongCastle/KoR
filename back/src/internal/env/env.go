package env

import (
	"log"
	"os"
)

var (
	envKeys = []string{
		"FRONT_HOST",      // フロントエンドのホスト
		"PASSWORD_PEPPER", // パスワードのペッパー
		"JWT_SECRET",      // JWT の秘密鍵
		"KOR_CERT",        // アプリケーションの cert
		"MYSQL_HOST",      // MySQL のホスト
		"MYSQL_PORT",      // MySQL のポート
		"MYSQL_DATABASE",  // MySQL のデータベース名
		"MYSQL_USER",      // MySQL のユーザ名
		"MYSQL_ROOT_PASSWORD_FILE", // MySQL ユーザのパスワード
	}
 	_env appEnv
)

type appEnv struct {
	Store map[string]string
}

func (e *appEnv) Init() {
	e.Store = make(map[string]string, len(envKeys))
	for _, key := range envKeys {
		e.Store[key] = os.Getenv(key)
	}
}

func init() {
	_env = appEnv{}
	_env.Init()
}

func Get(key string) string {
	v := _env.Store[key]
	if v == "" {
		log.Printf("[WARNING] Access to unknown env: %s\n", key)
	}
	return v
}
