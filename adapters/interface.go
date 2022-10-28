package adapters

import "fmt"

type Cacheable interface {
	fmt.Stringer
	GetKey() string
	GetRedisPair() (string, string, error) //key, value, error
	ParseValue(string) error
	Assign(Cacheable) error
}

type Cacher interface {
	Cache(cacheable Cacheable) error
	GetCachable(cacheable Cacheable) error 
}