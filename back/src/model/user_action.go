package model

import (
	"errors"

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
	// TODO: 検証する (引数含む)
	user, err := GetUser(&UserGetQuery{ID: userParams.ID})
	if err != nil {
		return nil, err
	}
	user.BindParams(userParams)

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
