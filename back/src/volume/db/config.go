package db

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

var myConf *MySQLConf

const PasswordMaxLen = 1024

// Only MySQL
func SetUp() error {
	if err := InitConf(); err != nil {
		return err
	}
	invoke()
	return nil
}

func InitConf() error {
	myConf = &MySQLConf{}
	return myConf.Init()
}

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
	if err := c.setPassword(); err != nil {
		return err
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
		c.userName, c.password, c.host, c.port, c.dbName,
	)
}

func (c *MySQLConf) setPassword() error {
	passwordPath := os.Getenv("MYSQL_ROOT_PASSWORD_FILE")
	if passwordPath == "" {
		return errors.New("Blank MYSQL_ROOT_PASSWORD_FILE")
	}

	f, err := os.Open(passwordPath)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, PasswordMaxLen)
	n, err := f.Read(buf)
	if err != nil {
		return err
	}
	pw := bytes.SplitN(buf, []byte("\n"), 2)[0]

	n, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	if n > 0 {
		return errors.New("Too Long DB Password")
	}

	c.password = string(pw)
	if c.password == "" {
		return errors.New("Empty DB Password")
	}
	return nil
}
