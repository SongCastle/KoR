package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

func newConnection() *Connection {
	return &Connection{
		client: newClient(),
		idledAt: time.Now(),
		Available: true,
	}
}

type Connection struct {
	client    *MySQLClient
	idledAt   time.Time
	Available bool
}

func (c *Connection) Lock() {
	c.Available = false
}

func (c *Connection) Unlock() {
	c.Available = true
	c.idledAt = time.Now()
}

func (c *Connection) Active() bool {
	return !c.Available || c.idledAt.Add(MaxIdleTime).After(time.Now())
}

func (c *Connection) Fresh() error {
	c.Lock()
	return c.client.Close()
}

func (c *Connection) DB() (*gorm.DB, error) {
	return c.client.Connect()
}
