package model

import (
	"errors"
	"reflect"
	"time"

	"github.com/SongCastle/KoR/db"
	"github.com/SongCastle/KoR/lib/encryptor"
	"github.com/SongCastle/KoR/lib/random"
	"github.com/jinzhu/gorm"
)

// password_salt の長さ
const UserPasswordSaltLen = 32

// Hooks
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Login == "" {
		return errors.New("Blank Login")
	}
	if u.Password == "" {
		return errors.New("Blank Password")
	}
	u.EncryptPassword()
	return nil
}

type User struct {
	ID                uint64 `json:"id,omitempty",gorm:"primaryKey"`
	Login             string `json:"login,omitempty"`
	Password          string `json:"password,omitempty"`
	Email             string `json:"email,omitempty"`
	PasswordSalt      string `json:"password_salt,omitempty"`
	EncryptedPassword string `json:"encrypted_password,omitempty"`
	AuthUUID          string `json:"auth_uuid,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

type UserGetQuery struct {
	ID    uint64
	Login string
}

func (u *User) EncryptPassword() error {
	if u.Password == "" {
		return nil
	}
	if u.PasswordSalt == "" {
		// Set PasswordSalt
		u.PasswordSalt = random.Generate(UserPasswordSaltLen)
	}
	// Encrypt Password
	digest, err := encryptor.Digest(u.Password + u.PasswordSalt)
	if err != nil {
		return err
	}
	u.EncryptedPassword = digest
	return nil
}

func (u *User) ValidPassword(password string) bool {
	return encryptor.Compare(u.EncryptedPassword, password + u.PasswordSalt)
}

// TODO: User オブジェクト内に関数群を含めた方がよいかも ...

func GetUsers(cols ...[]string) ([]User, error) {
	var users []User

	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()

	cdb := conn.DB()
	if len(cols) > 0 {
		cdb = cdb.Select(cols[0])
	}
	if err := cdb.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(query *UserGetQuery, cols ...[]string) (*User, error) {
	user := User{}
	user.ID, user.Login = query.ID, query.Login

	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()

	cdb := conn.DB().Where(&user)
	if len(cols) > 0 {
		cdb = cdb.Select(cols[0])
	}
	if err := cdb.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// TODO: 定義箇所の検討
type UserParams struct {
	ID       uint64
	Login    *string `json:"login"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	AuthUUID *string `json:"auth_uuid"`
}

func CreateUser(userParams *UserParams) (*User, error) {
	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()
	// TODO: 検証
	user := &User{}
	bindParamsToUser(userParams, user)
	user.ID, user.AuthUUID = 0, ""
	if err := conn.DB().Omit("Password").Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(userParams *UserParams) (*User, error) {
	// TODO: 検証する (引数含む)
	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()
	// TODO: DB Connection を Context で管理したい
	user, err := GetUser(&UserGetQuery{ID: userParams.ID})
	if err != nil {
		return nil, err
	}
	bindParamsToUser(userParams, user)
	if err := user.EncryptPassword(); err != nil {
		return nil, err
	}
	if err := conn.DB().Omit("Password").Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteUser(id uint64) error {
	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return err
	}
	defer conn.Close()
	// TODO: 検証
	if err := conn.DB().Delete(&User{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

func bindParamsToUser(params *UserParams, user *User) {
	rp, ru := reflect.ValueOf(params).Elem(), reflect.ValueOf(user).Elem()
	rpt := rp.Type()
	for i := 0; i < rpt.NumField(); i++ {
		// params のフィールド名を取得
		fn := rpt.Field(i).Name
		if v := rp.FieldByName(fn); !v.IsZero() {
			// user に同じフィールドがあるか確認
			if v2 := ru.FieldByName(fn); v2 != (reflect.Value{}) {
				if v.Kind() == v2.Kind() {
					// 同じ型のフィールドが存在する場合、値をセットする
					v2.Set(v)
				} else {
					// 型が違う場合、ポインタの参照先を確認する
					if iv := reflect.Indirect(v); iv.IsValid() {
						// 参照先の型が同じ型である場合、値をセットする
						if iv.Kind() == v2.Kind() {
							v2.Set(iv)
						}
					}
				}
			}
		}
	}
}
