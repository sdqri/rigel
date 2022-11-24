package adapters

import "fmt"

type Cacheable interface {
	fmt.Stringer
	GetKey(string) string
	GetPair(string) (string, string, error) //key, value, error
	ParseValue(value string) error
}

type Cacher interface {
	Cache(cacheable Cacheable) error
	GetCachable(cacheable Cacheable) error
}
