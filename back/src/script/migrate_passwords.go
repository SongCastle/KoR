// 暗号化していないパスワードを、暗号化したものへ移行します。
// $ MYSQL_PASSWORD=<password> go run script/migrate_passwords.go
// ### 注意点 ###
// 実行にあたり、環境変数 `PASSWORD_PEPPER` を設定してください。
// DB マイグレーションは 20220129012448 まで実施している必要があります。

package main

import (
	"fmt"
	"os"

	"github.com/SongCastle/KoR/db"
	"github.com/SongCastle/KoR/lib/encryptor"
	"github.com/SongCastle/KoR/model"
)

func main() {
	// DB 初期化
	if err := db.InitConf(); err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Failed to run.")
		return
	}
	// Encryptor 初期化
	encryptor.Init()

	// デフォルトパスワード設定
	defaultPassword := "password1234"
	if defaultPasswordEnv := os.Getenv("DEFAULT_USER_PASSWORD"); defaultPasswordEnv != "" {
		defaultPassword = defaultPasswordEnv
	}

	// DB 接続
	client := db.MySQLClient{}
	d, err := client.Connect()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Failed to connect DB.")
	}
	defer client.Close()

	// User 一覧取得
	var users []model.User
	if err := d.Select("id, password").Find(&users).Error; err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Failed to get Users list.")
		return
	}

	// パスワード移行
	for _, user := range users {
		// 1. password_salt を生成する
		// 2. password (生) を encrypted_password (暗号化) へ変更する

		fmt.Printf("User#%d\n", user.ID)
		fmt.Printf("Before Password: %s\n", user.Password)

		if user.Password == "" {
			user.Password = defaultPassword
		}
		if err := user.EncryptPassword(); err != nil {
			fmt.Printf("Failed to encrypt Password (User#%d)\n\n", user.ID)
			continue
		}

		if err := d.Model(user).Update(user).Error; err != nil {
			fmt.Printf("Failed to Update (User#%d)\n", user.ID)
			fmt.Printf("Error: %v\n\n", err)
		} else {
			fmt.Printf("Success to Update (User#%d)\n", user.ID)
			fmt.Printf("Encrypted Password: %s\n", user.EncryptedPassword)
			fmt.Printf("Password Validity: %t\n\n", user.ValidPassword(user.Password))
		}
	}
}
