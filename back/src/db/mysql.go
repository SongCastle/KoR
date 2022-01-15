package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var myConf *MySQLConf

type MySQLConf struct {
	driver   string
	host     string
	port     string
	dbName   string
	userName string
	password string
}

func (c *MySQLConf) Init() error {
	if c.host = os.Getenv("MYSQL_HOST"); c.host == "" {
		return errors.New("Blank MYSQL_HOST")
	}
	if c.dbName = os.Getenv("MYSQL_DATABASE"); c.dbName == "" {
		return errors.New("Blank MYSQL_DATABASE")
	}
	if c.port = os.Getenv("MYSQL_PORT"); c.port == "" {
		return errors.New("Blank MYSQL_PORT")
	}
	if c.userName = os.Getenv("MYSQL_USERNAME"); c.userName == "" {
		return errors.New("Blank MYSQL_USERNAME")
	}
	if c.password = os.Getenv("MYSQL_PASSWORD"); c.password == "" {
		return errors.New("Blank MYSQL_PASSWORD")
	}
	c.driver = "mysql"
	return nil
}

func (c *MySQLConf) Driver() string {
	return c.driver
}

func (c *MySQLConf) URL() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.userName,	c.password, c.host, c.port, c.dbName,
	)
}

type MySQLConn struct {
	dialect string
	url     string
	db      *gorm.DB
}

func (conn *MySQLConn) Open() error {
	var err error
	conn.db, err = gorm.Open(conn.dialect, conn.url)
	if err != nil {
		return err
	}
	return nil
}

func (conn *MySQLConn) Close() error {
	if conn.db == nil {
		return errors.New("Not Open")
	}
	return conn.db.Close()
}

func (conn *MySQLConn) DB() *gorm.DB {
	return conn.db
}

func initMySQL() error {
	myConf = &MySQLConf{}
	if err := myConf.Init(); err != nil {
		return err
	}
	return nil
}

func newMySQL() *MySQLConn {
	return &MySQLConn{dialect: myConf.Driver(), url: myConf.URL()}
}
