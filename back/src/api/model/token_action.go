package model

import (
	"errors"

	"github.com/SongCastle/KoR/volume/db"
	"github.com/jinzhu/gorm"
)

// Hooks
func (t *Token) BeforeCreate(tx *gorm.DB) error {
	if t.UserID == 0 {
		return errors.New("Empty UserID")
	}
	if t.UUID == "" {
		return errors.New("Empty UUID")
	}
	return nil
}

type TokenGetQuery struct {
	UserID uint64
	UUID   string
}

type TokenParams struct {
	ID     uint64
	UserID uint64  `json:"user_id"`
	UUID   *string `json:"uuid"`
}

func GetToken(query *TokenGetQuery) (*Token, error) {
	token := &Token{UserID: query.UserID, UUID: query.UUID}

	err := db.Connect(func(d *gorm.DB) error {
		return d.Where(token).First(token).Error
	})

	if err == nil {
		return token, nil
	}
	return nil, err
}

func CreateToken(params *TokenParams) (*Token, error) {
	token := &Token{}
	token.BindParams(params)
	token.ID = 0

	err := db.Connect(func(d *gorm.DB) error {
		return d.Create(token).Error
	})

	if err == nil {
		return token, nil
	}
	return nil, err
}
