package db

import "github.com/jinzhu/gorm"

type DBConn interface {
	Open() error
	Close() error
	DB() *gorm.DB
}

// Only MySQL
func InitDB() error {
	return initMySQL()
}

func NewDB() DBConn {
	return newMySQL()
}
