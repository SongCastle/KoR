package model

import (
	"errors"

	"github.com/SongCastle/KoR/internal/random"
	"github.com/SongCastle/KoR/volume/db"
	"github.com/jinzhu/gorm"
)

const UUIDLen = 32

// Hooks
func (t *Token) BeforeCreate(tx *gorm.DB) error {
	if t.UserID == 0 {
		return errors.New("Empty UserID")
	}
	if t.UUID == "" {
		t.UUID = random.Generate(UUIDLen)
	}
	return nil
}

func (t *Token) BeforeDelete(tx *gorm.DB) error {
	authes, err := GetAuthorities(func (db *gorm.DB) *gorm.DB {
		return db.Select([]string{"id"}).Where(&Authority{TokenID: t.ID})
	})
	if err != nil {
		return err
	}
	for _, auth := range authes {
		err := DeleteAuthority(auth.ID)
		if err != nil {
			return err
		}
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

func GetTokens(query func(*gorm.DB) *gorm.DB) ([]Token, error) {
	var tokens []Token
	err := db.Connect(func(d *gorm.DB) error {
		return query(d).Find(&tokens).Error
	})
	if err == nil {
		return tokens, nil
	}
	return nil, err
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

func DeleteToken(id uint64) error {
	return db.Connect(func(d *gorm.DB) error {
		return d.Delete(&Token{ID: id}).Error
	})
}
