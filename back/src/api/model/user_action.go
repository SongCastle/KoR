package model

import (
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

func (u *User) BeforeDelete(tx *gorm.DB) error {
	tokens, err := GetTokens(func (db *gorm.DB) *gorm.DB {
		return db.Select([]string{"id"}).Where(&Token{UserID: u.ID})
	})
	if err != nil {
		return err
	}
	for _, token := range tokens {
		err := DeleteToken(token.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

type UserParams struct {
	ID       uint64
	Login    *string `json:"login"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
}

type UserGetQuery struct {
	ID    uint64
	Login string
}

func GetUsers(cols ...string) ([]User, error) {
	var users []User

	err := db.Connect(func(d *gorm.DB) error {
		var err error
		if len(cols) > 0 {
			err = d.Select(cols).Find(&users).Error
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
