package adapters

import "fmt"

type Pair[T any] struct {
	Key   string
	Value T
}

type Cacheable[T any] interface {
	fmt.Stringer
	GetKey() (string, error)
	GetPair() (Pair[T], error) //pair{key, value}, error
	ParseValue(value T) error
}

type Cacher[T any] interface {
	Cache(cacheable Cacheable[T]) error
	GetCachable(cacheable Cacheable[T]) error
}

type MultilevelCacher[T any] struct {
	Cachers []Cacher[T]
}

func NewMultilevelCacher[T any](cachers ...Cacher[T]) *MultilevelCacher[T] {
	return &MultilevelCacher[T]{
		Cachers: cachers,
	}
}

func (mlc *MultilevelCacher[T]) Cache(cacheable Cacheable[T]) (err error) {
	for _, cacher := range mlc.Cachers {
		err = cacher.Cache(cacheable)
		if err != nil {
			return
		}
	}
	return nil
}

func (mlc *MultilevelCacher[T]) GetCachable(cacheable Cacheable[T]) (err error) {
	for _, cacher := range mlc.Cachers {
		err = cacher.GetCachable(cacheable)
		if err == nil {
			return
		}
	}
	return
}
