package adapters

import (
	"encoding/json"
	"errors"

	"github.com/Code-Hex/go-generics-cache/policy/lfu"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCacheableType error = errors.New("cacheable should be of type lfuCache")
	ErrValueType     error = errors.New("cached value type is strange")
)

type MemoryClient[T Cacheable] struct {
	LogEntry *log.Entry
	Prefix   string
	lfuCache *lfu.Cache[string, T]
}

func NewMemoryClient[T Cacheable](logEntry *log.Entry, prefix string, cap int) *MemoryClient[T] {
	// Setting package specific fields for log entry
	entry := logEntry.WithFields(log.Fields{
		"package": "adapters.memory",
	})
	c := lfu.NewCache[string, T](lfu.WithCapacity(cap))
	return &MemoryClient[T]{
		LogEntry: entry,
		Prefix:   prefix,
		lfuCache: c,
	}
}

func (mc *MemoryClient[T]) Cache(cacheable Cacheable) error {
	cacheableT, ok := cacheable.(T)
	if !ok {
		return ErrCacheableType
	}

	mc.lfuCache.Set(cacheable.GetKey(mc.Prefix), cacheableT)
	return nil
}

func (mc *MemoryClient[T]) GetCachable(cacheable Cacheable) error {
	memCacheable, ok := mc.lfuCache.Get(cacheable.GetKey(mc.Prefix))
	if !ok {
		return ErrKeyNotExists
	}

	value, err := json.Marshal(memCacheable)
	if err != nil {
		return ErrValueType
	}
	err = cacheable.ParseValue(string(value))
	if err != nil {
		return err
	}
	return nil
}
