package model

import "time"

type Token struct {
	ID          uint64      `json:"id,omitempty",gorm:"primaryKey"`
	UserID      uint64      `json:"user_id,omitempty"`
	UUID        string      `json:"uuid,omitempty"`
	User        User        `json:"-",gorm:"foreignKey:UserID"`
	Authorities []Authority `json:"-",gorm:"foreignKey:TokenID"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
}

func (t *Token) BindParams(params *TokenParams) {
	bindParamsToObject(params, t)
}
