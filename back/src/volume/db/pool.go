package db

import (
	"errors"
	"sync"
	"time"

	"github.com/SongCastle/KoR/internal/log"
	"github.com/jinzhu/gorm"
)

const (
	PoolSize    = 5
	RetryTimes  = 5
	WaitTime    = time.Second
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
						log.Infof("Reset Connection %d (err: %v)", i, c.Fresh())
						pool[i] = nil
						mu.Unlock()
					}
				}
				time.Sleep(MaxIdleTime)
			}
		}
	}()
}

func retrieveConnection() (*Connection, error) {
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

func connection() (*Connection, error) {
	c, err := retrieveConnection()
	if err == nil {
		return c, nil
	}
	for i := 0; i < RetryTimes; i++ {
		time.Sleep(WaitTime)
		c, err = retrieveConnection()
		if err == nil {
			return c, nil
		}
	}
	return nil, err
}

func Connect(f func(db *gorm.DB) error) error {
	c, err := connection()
	if err != nil {
		return err
	}
	defer c.Unlock()

	db, err := c.DB()
	if err != nil {
		return err
	}
	// TODO: Rollback について
	return db.Transaction(f)
}
