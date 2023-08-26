package adapters

import (
	"fmt"

	"github.com/Code-Hex/go-generics-cache/policy/lru"
	log "github.com/sirupsen/logrus"
)

type MemoryClient struct {
	LogEntry *log.Entry
	Prefix   string
	lfuCache *lru.Cache[string, string]
}

func NewMemoryClient(logEntry *log.Entry, prefix string, cap int) *MemoryClient {
	// Setting package specific fields for log entry
	entry := logEntry.WithFields(log.Fields{
		"package": "adapters.memory",
	})
	c := lru.NewCache[string, string](lru.WithCapacity(cap))
	return &MemoryClient{
		LogEntry: entry,
		Prefix:   prefix,
		lfuCache: c,
	}
}

func (mc *MemoryClient) Cache(cacheable Cacheable[string]) error {
	mc.LogEntry.Debugf("writing cacheable(%s) to memory(lru cache)", cacheable.String())
	pair, err := cacheable.GetPair()
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%s", mc.Prefix, pair.Key)
	value := pair.Value
	mc.lfuCache.Set(key, value)
	return nil
}

func (mc *MemoryClient) GetCachable(cacheable Cacheable[string]) error {
	key, err := cacheable.GetKey()
	if err != nil {
		return err
	}
	key = fmt.Sprintf("%s:%s", mc.Prefix, key)
	mc.LogEntry.Debugf("reading cacheable(%s) from memory(lru cache)", key)
	value, ok := mc.lfuCache.Get(key)
	if !ok {
		return ErrKeyNotExists
	}

	err = cacheable.ParseValue(value)
	if err != nil {
		return err
	}
	return nil
}
