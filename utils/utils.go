package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

func SHA1Sum(text string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(text))
	return hex.EncodeToString(sha1.Sum(nil))
}

func mapKeysToSlice[T comparable, V any](m map[T]V) []T {
	keys := make([]T, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func sign(key string, salt string, input []byte) string {
	h := hmac.New(sha1.New, []byte(key))
	input = append(input, []byte(salt)...)
	h.Write(input)
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
