package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
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

type Signatory struct {
	prefix string
	key    string
	salt   string
	r      *regexp.Regexp
}

func NewSignatory(prefix string, key string, salt string) *Signatory {
	r := regexp.MustCompile(fmt.Sprintf(`%s/(.+)`, prefix))
	return &Signatory{prefix: prefix, key: key, salt: salt, r: r}
}

func (signatory *Signatory) SignURL(u url.URL) string {
	SignableSlice := make([]string, 0)

	// Extract & Add requets_path to SignableSlice
	groups := signatory.r.FindStringSubmatch(u.Path)
	var requestPath string
	if len(groups) >= 2 {
		requestPath = groups[1]
	}
	SignableSlice = append(SignableSlice, fmt.Sprintf("%s=%s", "request_path", requestPath))

	// Note: singature is unique for query fields with list of values because only the first element is used in checking (line 69)
	// Extract & Add sorted QueryParams to SignableSlice
	queryValues := u.Query()
	queryKeys := mapKeysToSlice(queryValues)
	sort.Strings(queryKeys)
	for _, k := range queryKeys {
		if k != "X-Signature" {
			queryPair := fmt.Sprintf("%s=%s", k, queryValues[k][0])
			SignableSlice = append(SignableSlice, queryPair)
		}
	}

	SignableBytes := []byte(strings.Join(SignableSlice, "&"))

	return sign(signatory.key, signatory.salt, SignableBytes)
}
