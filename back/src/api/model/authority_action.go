package model

import (
	"fmt"

	"github.com/SongCastle/KoR/volume/db"
	"github.com/jinzhu/gorm"
)

// Hooks
func (a *Authority) BeforeCreate(tx *gorm.DB) error {
	if a.TokenID == 0 {
		return fmt.Errorf("Empty TokenID")
	}
	if a.Type == "" {
		return fmt.Errorf("Empty Type")
	}
	if a.Right == nil {
		return fmt.Errorf("Empty Right")
	}
	if len(a.Right) != 8 {
		return fmt.Errorf("Invalid Right (len: %d)", len(a.Right))
	}
	return nil
}

type AuthorityParams struct {
	ID        uint64
	TokenID   uint64  `json:"token_id"`
	Type      *string  `json:"type"`
	Right     []uint8 `json:"right"`
}

func GetAuthority(id uint64) (*Authority, error) {
	auth := &Authority{ID: id}
	err := db.Connect(func(d *gorm.DB) error {
		return d.Where(auth).First(auth).Error
	})
	if err == nil {
		return auth, nil
	}
	return nil, err
}

func CreateAuthority(params *AuthorityParams) (*Authority, error) {
	auth := &Authority{}
	auth.BindParams(params)
	auth.ID = 0

	err := db.Connect(func(d *gorm.DB) error {
		return d.Create(&auth).Error
	})

	if err == nil {
		return auth, nil
	}
	return nil, err
}
