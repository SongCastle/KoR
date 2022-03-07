package model

import (
	"reflect"
	"time"

	"github.com/SongCastle/KoR/lib/encryptor"
	"github.com/SongCastle/KoR/lib/random"
)

// password_salt の長さ
const UserPasswordSaltLen = 32

type UserParams struct {
	ID       uint64
	Login    *string `json:"login"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	AuthUUID *string `json:"-"`
}

type User struct {
	ID                uint64 `json:"id,omitempty",gorm:"primaryKey"`
	Login             string `json:"login,omitempty"`
	Password          string `json:"password,omitempty"`
	Email             string `json:"email,omitempty"`
	PasswordSalt      string `json:"-,omitempty"`
	EncryptedPassword string `json:"-,omitempty"`
	AuthUUID          string `json:"-,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
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
	u.Password = ""
	return nil
}

func (u *User) ValidPassword(password string) bool {
	return encryptor.Compare(u.EncryptedPassword, password + u.PasswordSalt)
}

func (u *User) BindParams(params *UserParams) {
	bindParamsToObject(params, u)
}

func bindParamsToObject(params interface{}, obj interface{}) {
	vp := reflect.ValueOf(params)
	ivp := reflect.Indirect(vp)
	// 引数が nil ポインタである、または非参照型である (Elem 取得不可)
	if !ivp.IsValid() || vp == ivp {
		return
	}
	vo := reflect.ValueOf(obj)
	ivo := reflect.Indirect(vo)
	if !ivo.IsValid() || vo == ivo {
		return
	}
	rp, ro := vp.Elem(), vo.Elem()
	rpt := rp.Type()
	for i := 0; i < rpt.NumField(); i++ {
		// params のフィールド名を取得
		fn := rpt.Field(i).Name
		if v := rp.FieldByName(fn); !v.IsZero() {
			// obj に同じフィールドがあるか確認
			if v2 := ro.FieldByName(fn); v2 != (reflect.Value{}) {
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
