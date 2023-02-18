package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type RedisClient struct {
	LogEntry   *log.Entry
	Prefix     string
	Client     *redis.Client
	Expiration time.Duration
}

func NewRedisClient(logEntry *log.Entry, prefix string, addr string, password string, db int, timeout time.Duration, expiration time.Duration) *RedisClient {
	// Setting package specific fields for log entry
	entry := logEntry.WithFields(log.Fields{
		"package": "adapters.redis",
	})

	// Creating RedisClient
	rc := RedisClient{
		LogEntry: entry,
		Prefix:   prefix,
		Client: redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     password,
			DB:           db,
			DialTimeout:  timeout,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		}),
		Expiration: expiration,
	}
	return &rc
}

func (rc *RedisClient) Cache(cacheable Cacheable[string]) error {
	rc.LogEntry.Debugf("writing cacheable(%s) to redis", cacheable.String())
	ctx := context.Background()
	pair, err := cacheable.GetPair()
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%s", rc.Prefix, pair.Key)
	value := pair.Value
	return rc.Client.Set(ctx, key, value, rc.Expiration).Err()
}

func (rc *RedisClient) GetCachable(cacheable Cacheable[string]) error {
	ctx := context.Background()
	key, err := cacheable.GetKey()
	if err != nil {
		return err
	}
	key = fmt.Sprintf("%s:%s", rc.Prefix, key)
	rc.LogEntry.Debugf("reading cacheable(%s) from redis", key)
	value, err := rc.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrKeyNotExists
	}

	err = cacheable.ParseValue(value)
	if err != nil {
		return err
	}
	return nil
}
