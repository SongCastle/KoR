package model

import (
	"errors"
	"time"
)

type Token struct {
	ID          uint64      `json:"id,omitempty",gorm:"primaryKey"`
	UserID      uint64      `json:"user_id,omitempty"`
	UUID        string      `json:"uuid,omitempty"`
	User        User        `json:"-",gorm:"foreignKey:UserID"`
	Authorities []Authority `json:"-",gorm:"foreignKey:TokenID"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
}

func (t *Token) IsPersisted() bool {
	return t.ID != 0
}

func (t *Token) AddUserAuthority(all bool) error {
	if !t.IsPersisted() {
		return errors.New("Not Persisted")
	}
	_, err := CreateAuthority(
		WithTokenID(t.ID),
		WithUsersRight(),
		WithCreateRight(all),
		WithReadRight(all),
		WithUpdateRight(all),
		WithDeleteRight(all),
	)
	return err
}

func (t *Token) BindParams(params *TokenParams) {
	bindParamsToObject(params, t)
}
