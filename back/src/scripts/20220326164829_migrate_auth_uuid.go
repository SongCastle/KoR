//
// users テーブルの auth_uuid カラムを tokens テーブルの uuid へ移行します。
// 並行して、 token の権限付与も実施します。(ユーザの CRUD 権限)
//
// ```
// $ go run scripts/20220326164829_migrate_auth_uuid.go
// ```
// ### 注意点 ###
// 実行にあたり、DB のマイグレーションを 20220326164829 まで実施している必要があります。
//

package main

import (
	"fmt"

	"github.com/SongCastle/KoR/api/model"
	"github.com/SongCastle/KoR/volume/db"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID       uint64 `json:"id,omitempty",gorm:"primaryKey"`
	Login    string `json:"login,omitempty"`
	AuthUUID string `json:"-,omitempty"`
}

func main() {
	if err := db.InitConf(); err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Failed to run.")
		return
	}

	var users []User
	const all = false

	err := db.Connect(func(d *gorm.DB) error {
		return d.Select([]string{"id", "login", "auth_uuid"}).Find(&users).Error
	})
	if err != nil {
		fmt.Printf("Users Error: %v\n", err)
		return
	}

	fmt.Println("*** Start Migration ***")

	for _, user := range users {
		if user.AuthUUID == "" {
			fmt.Printf("- User / id: %v, uuid: %v ... SKIP\n", user.ID, user.AuthUUID)
			continue
		}
		fmt.Printf("- User / id: %v, uuid: %v\n", user.ID, user.AuthUUID)

		tParams := &model.TokenParams{
			UserID: user.ID, UUID: &user.AuthUUID,
		}
		token := model.NewToken(tParams)
		if err := token.Create(); err != nil {
			fmt.Printf("Failed to Token: %v\n", err)
			continue
		}

		aParams := &model.AuthorityParams{TokenID: token.ID}
		auth := model.NewAuthority(aParams)
		auth.
			AsUsersRight().
			AddCreateRight(all).
			AddReadRight(all).
			AddUpdateRight(all).
			AddDeleteRight(all)

		if err := auth.Create(); err != nil {
			fmt.Printf("Failed to Create Authority: %v\n", err)
			continue
		}

		fmt.Printf(" - New Token / id: %d, uuid: %s\n", token.ID, token.UUID)
		fmt.Printf("  - New Authority / id: %d, create: %t, read: %t, update: %t, delete: %t\n\n",
			auth.ID, auth.CanCreate(all), auth.CanRead(all), auth.CanUpdate(all), auth.CanDelete(all),
		)
	}
	fmt.Printf("Finish Updating AuthUUID\n\n")

	fmt.Println("*** Check Migration ***")

	for _, user := range users {
		if user.AuthUUID == "" {
			fmt.Printf("- User / id: %v, uuid: %v ... SKIP\n", user.ID, user.AuthUUID)
			continue
		}
		fmt.Printf("- User / id: %v, uuid: %v\n", user.ID, user.AuthUUID)

		tParams := &model.TokenParams{UserID: user.ID, UUID: &user.AuthUUID}
		token, err := model.GetToken(model.WhereToken(tParams))
		if err != nil {
			fmt.Printf("Failed to Get Token: %v\n", err)
			continue
		}

		aParams := &model.AuthorityParams{TokenID: token.ID}
		auth, err := model.GetAuthority(model.WhereAuthority(aParams))
		if err != nil {
			fmt.Printf("Failed to Get Authority: %v\n", err)
			continue
		}

		fmt.Printf(" - Token / id: %d, uuid: %s\n", token.ID, token.UUID)
		fmt.Printf("  - Authority / id: %d, create: %t, read: %t, update: %t, delete: %t\n\n",
			auth.ID, auth.CanCreate(all), auth.CanRead(all), auth.CanUpdate(all), auth.CanDelete(all),
		)
	}
	fmt.Println("Finish Migration")
}
