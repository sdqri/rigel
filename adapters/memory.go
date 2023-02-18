package adapters

import (
	"fmt"

	"github.com/Code-Hex/go-generics-cache/policy/lfu"
	log "github.com/sirupsen/logrus"
)

type MemoryClient struct {
	LogEntry *log.Entry
	Prefix   string
	lfuCache *lfu.Cache[string, string]
}

func NewMemoryClient(logEntry *log.Entry, prefix string, cap int) *MemoryClient {
	// Setting package specific fields for log entry
	entry := logEntry.WithFields(log.Fields{
		"package": "adapters.memory",
	})
	c := lfu.NewCache[string, string](lfu.WithCapacity(cap))
	return &MemoryClient{
		LogEntry: entry,
		Prefix:   prefix,
		lfuCache: c,
	}
}

func (mc *MemoryClient) Cache(cacheable Cacheable[string]) error {
	mc.LogEntry.Debugf("writing cacheable(%s) to memory(lfu cache)", cacheable.String())
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
	mc.LogEntry.Debugf("reading cacheable(%s) from memory(lfu cache)", key)
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
