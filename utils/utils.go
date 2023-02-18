package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func SHA1Sum(text string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(text))
	return hex.EncodeToString(sha1.Sum(nil))
}

func Sign(key string, salt string, input []byte) string {
	h := hmac.New(sha1.New, []byte(key))
	input = append(input, []byte(salt)...)
	h.Write(input)
	return hex.EncodeToString(h.Sum(nil))
}
