package adapters

import (
	"errors"

	"github.com/Code-Hex/go-generics-cache/policy/lfu"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCacheableType error = errors.New("cacheable should be of type lfuCache")
)

type MemoryClient[T Cacheable] struct {
	LogEntry            *log.Entry
	redisClient *RedisClient
	lfuCache *lfu.Cache[string, T]
}

func NewMemoryClient[T Cacheable](logEntry *log.Entry, redisClient *RedisClient, cap int) *MemoryClient[T] {
	// Setting package specific fields for log entry
	entry := logEntry.WithFields(log.Fields{
		"package": "adapters.memory",
	})
	c := lfu.NewCache[string, T](lfu.WithCapacity(cap))
	return &MemoryClient[T]{
		LogEntry: entry,
		redisClient: redisClient,
		lfuCache: c,
	}
}

func (mc *MemoryClient[T]) Cache(cacheable Cacheable) error {
	err := mc.redisClient.Cache(cacheable)
	if err != nil{
		return err
	}

	cacheableT, ok := cacheable.(T)
	if !ok{
		return ErrCacheableType
	}

	mc.lfuCache.Set(cacheable.GetKey(), cacheableT)
	return nil
}

func (mc *MemoryClient[T]) GetCachable(cacheable Cacheable) error {
	memCacheable, ok := mc.lfuCache.Get(cacheable.GetKey())
	if !ok{
		return mc.redisClient.GetCachable(cacheable)
	}
	cacheable.Assign(memCacheable)
	return nil
}
