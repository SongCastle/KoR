package model

import (
	"errors"
	"reflect"

	"github.com/SongCastle/KoR/db"
	"github.com/jinzhu/gorm"
)

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

type UserGetQuery struct {
	ID    uint64
	Login string
}

type UserParams struct {
	ID       uint64
	Login    *string `json:"login"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	AuthUUID *string `json:"-"`
}

func GetUsers(cols ...[]string) ([]User, error) {
	var users []User

	err := db.Connect(func(d *gorm.DB) error {
		var err error
		if len(cols) > 0 {
			err = d.Select(cols[0]).Find(&users).Error
		} else {
			err = d.Find(&users).Error
		}
		return err
	})

	if err == nil {
		return users, nil
	}
	return nil, err
}

func GetUser(query *UserGetQuery, cols ...[]string) (*User, error) {
	user := User{ID: query.ID, Login: query.Login}

	err := db.Connect(func(d *gorm.DB) error {
		var err error
		if len(cols) > 0 {
			err = d.Where(&user).Select(cols[0]).First(&user).Error
		} else {
			err = d.Where(&user).First(&user).Error
		}
		return err
	})

	if err == nil {
		return &user, nil
	}
	return nil, err
}

func CreateUser(userParams *UserParams) (*User, error) {
	// TODO: 検証
	user := &User{}
	bindParamsToUser(userParams, user)
	user.ID = 0

	err := db.Connect(func(d *gorm.DB) error {
		return d.Omit("Password").Create(user).Error
	})

	if err == nil {
		return user, nil
	}
	return nil, err
}

func UpdateUser(userParams *UserParams) (*User, error) {
	// TODO: 検証する (引数含む)
	user, err := GetUser(&UserGetQuery{ID: userParams.ID})
	if err != nil {
		return nil, err
	}
	bindParamsToUser(userParams, user)

	if err := user.EncryptPassword(); err != nil {
		return nil, err
	}

	err = db.Connect(func(d *gorm.DB) error {
		return d.Omit("Password").Save(user).Error
	})

	if err == nil {
		return user, nil
	}
	return nil, err
}

func DeleteUser(id uint64) error {
	// TODO: 検証, 存在しないユーザの対策
	return db.Connect(func(d *gorm.DB) error {
		return d.Delete(&User{ID: id}).Error
	})
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
