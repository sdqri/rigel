package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var(
	ErrKeyNotExists error = errors.New("key doesnt not exists")
)

type RedisClient struct {
	LogEntry            *log.Entry
	Client              *redis.Client
}

func NewRedisClient(logEntry *log.Entry, addr string, password string, db int, timeout time.Duration) *RedisClient {
	// Setting package specific fields for log entry
	entry := logEntry.WithFields(log.Fields{
		"package": "adapters.redis",
	})

	// Creating RedisClient
	rc := RedisClient{
		LogEntry: entry,
		Client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
			DialTimeout: timeout,
			ReadTimeout: timeout,
			WriteTimeout: timeout,
		}),
	}
	return &rc
}

func (rc *RedisClient) Cache(cacheable Cacheable) error {
	rc.LogEntry.Debugf("writing cacheable(%s) to redis", cacheable.String())
	ctx := context.Background()
	key, value, err := cacheable.GetRedisPair()
	if err!= nil{
		return err
	}
	return rc.Client.Set(ctx, key, value, 0).Err()
}

func (rc *RedisClient) GetCachable(cacheable Cacheable) error {
	ctx := context.Background()
	key := cacheable.GetKey()
	rc.LogEntry.Debugf("reading cacheable(%s) from redis", key)
	value, err := rc.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrKeyNotExists
	}
	
	err = cacheable.ParseValue(value)
	if err != nil{
		return err
	}
	return nil
}
