package tommorrowio

import lru "github.com/hashicorp/golang-lru"

func withCache[T any](cache *lru.Cache, key string, fn func() (T, error)) (T, error) {
	if cached, ok := cache.Get(key); ok {
		return cached.(T), nil
	}

	value, err := fn()
	if err != nil {
		return value, err
	}

	cache.Add(key, value)

	return value, nil
}
