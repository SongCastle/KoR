package model

import (
	"reflect"

	"github.com/SongCastle/KoR/volume/db"
	"github.com/jinzhu/gorm"
)

// Hooks
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.ValidateLogin(); err != nil {
		return err
	}
	if err := u.ValidatePassword(); err != nil {
		return err
	}
	if err := u.EncryptPassword(); err != nil {
		return err
	}
	return nil
}

type UserParams struct {
	ID       uint64
	Login    *string `json:"login"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	AuthUUID *string `json:"-"`
}

type UserGetQuery struct {
	ID    uint64
	Login string
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
	user := &User{}
	user.BindParams(userParams)
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
	user, err := GetUser(&UserGetQuery{ID: userParams.ID})
	if err != nil {
		return nil, err
	}
	user.BindParams(userParams)

	if userParams.Login != nil {
		if err := user.ValidateLogin(); err != nil {
			return nil, err
		}
	}
	if userParams.Password != nil {
		if err := user.ValidatePassword(); err != nil {
			return nil, err
		}
		if err := user.EncryptPassword(); err != nil {
			return nil, err
		}
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
	// TODO: 存在しないユーザの場合について
	return db.Connect(func(d *gorm.DB) error {
		return d.Delete(&User{ID: id}).Error
	})
}

func existUser(query *UserGetQuery) (bool, error) {
	count := int64(0)

	err := db.Connect(func(d *gorm.DB) error {
		user := &User{ID: query.ID, Login: query.Login}
		return d.Model(user).Where(user).Count(&count).Error
	})
	return count > 0, err
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
