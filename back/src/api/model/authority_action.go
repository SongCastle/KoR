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

func GetAuthorities(query func(*gorm.DB) *gorm.DB) ([]Authority, error) {
	var authes []Authority
	err := db.Connect(func(d *gorm.DB) error {
		return query(d).Find(&authes).Error
	})
	if err == nil {
		return authes, nil
	}
	return nil, err
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

func CreateAuthority(adds ...AuthorityBuilder) (*Authority, error) {
	auth := &Authority{}
	auth.Build(adds...)

	err := db.Connect(func(d *gorm.DB) error {
		return d.Create(&auth).Error
	})

	if err == nil {
		return auth, nil
	}
	return nil, err
}

func DeleteAuthority(id uint64) error {
	return db.Connect(func(d *gorm.DB) error {
		return d.Delete(&Authority{ID: id}).Error
	})
}
