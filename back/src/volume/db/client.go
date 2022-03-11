package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func newClient() *MySQLClient {
	client := MySQLClient{}
	client.Init()
	return &client
}

type MySQLClient struct {
	dialect string
	url     string
	db      *gorm.DB
}

func (c *MySQLClient) Init() {
	c.dialect, c.url = myConf.Driver(), myConf.URL()
}

func (c *MySQLClient) Connect() (*gorm.DB, error) {
	if c.db != nil {
		return c.db, nil
	}
	var err error
	c.db, err = gorm.Open(c.dialect, c.url)
	if err == nil {
		return c.db, nil
	}
	return nil, err
}

func (c *MySQLClient) Close() error {
	if c.db == nil {
		return nil
	}
	return c.db.Close()
}
