package gorigelsdk

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
)

func SerializeToQueryString(obj map[string]interface{}) string {
	var str []string
	for p, v := range obj {
		str = append(str, fmt.Sprintf("%s=%v", p, v))
	}
	return strings.Join(str, "&")
}

func Sign(key string, salt string, input string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(input + salt))
	b64Signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	b64Signature = strings.TrimRight(b64Signature, "=")
	return b64Signature
}

func SignQueryString(key string, salt string, requestPath string, queryString string, expiry int64) string {
	signableSlice := []string{}
	signableSlice = append(signableSlice, fmt.Sprintf("request_path=%s", requestPath))

	if queryString != "" {
		querySlice := strings.Split(queryString, "&")
		if expiry != 0 && expiry != -1 {
			querySlice = append(querySlice, fmt.Sprintf("X-ExpiresAt=%d", expiry))
		}
		sort.Strings(querySlice)
		signableSlice = append(signableSlice, querySlice...)
	}

	signableString := strings.Join(signableSlice, "&")
	signature := Sign(key, salt, signableString)
	signableSlice = append(signableSlice, fmt.Sprintf("X-Signature=%s", signature))
	signableSlice = signableSlice[1:]
	return strings.Join(signableSlice, "&")
}

func SerializeMapToQueryString(params map[string]string) string {
	var parts []string
	for key, value := range params {
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(parts, "&")
}

func Ptr[T any](v T) *T {
	return &v
}
