package db

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	PoolSize    = 5
	WaitTime    = 5 * time.Second
	MaxIdleTime = 3 * time.Minute
)

var (
	pool         [PoolSize]*Connection
	mu           sync.Mutex
	timeoutError error = errors.New("connection timeout")
	stopCleaner  chan bool = make(chan bool)
)

func invoke() {
	go func() {
		for {
			select {
			case <-stopCleaner: // TODO: Graceful Shutdown に合わせたい
				return
			default:
				for i, c := range pool {
					if !(c == nil || c.Active()) {
						mu.Lock()
						log.Printf("Reset Connection %d (err: %v)\n", i, c.Fresh())
						pool[i] = nil
						mu.Unlock()
					}
				}
				time.Sleep(MaxIdleTime)
			}
		}
	}()
}

func connection() (*Connection, error) {
	mu.Lock()
	defer mu.Unlock()
	for i, c := range pool {
		if c == nil {
			pool[i] = newConnection()
			pool[i].Lock()
			return pool[i], nil
		} else if c.Available {
			c.Lock()
			return c, nil
		}
	}
	return nil, timeoutError
}

func Connect(f func(db *gorm.DB) error) error {
	c, err := connection()
	if err != nil {
		if err != timeoutError {
			return err
		}
		time.Sleep(WaitTime)
		c, err = connection()
		if err != nil {
			return err
		}
	}
	defer c.Unlock()

	db, err := c.DB()
	if err != nil {
		return err
	}
	return db.Transaction(f)
}
