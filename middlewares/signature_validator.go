package middlewares

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdqri/rigel/utils"
)

func MapKeysToSlice[T comparable, V any](m map[T]V) []T {
	keys := make([]T, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

type XParams struct {
	Signature string `form:"X-Signature" binding:"required"`
	ExpiresAt int64  `form:"X-ExpiresAt"`
}

func NewSignatureValidator(key string, salt string, prefix string) gin.HandlerFunc {
	r := regexp.MustCompile(fmt.Sprintf(`%s/(.+)`, prefix))
	return func(c *gin.Context) {
		var args XParams
		err := c.Bind(&args)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		SignableSlice := make([]string, 0)

		groups := r.FindStringSubmatch(c.Request.URL.Path)
		var requestPath string
		if len(groups) >= 2 {
			requestPath = groups[1]
		}
		SignableSlice = append(SignableSlice, fmt.Sprintf("%s=%s", "request_path", requestPath))

		// Note: singature is uniquer for query fields with list of values because only the first element is used in checking (line 58)
		queryValues := c.Request.URL.Query()
		queryKeys := MapKeysToSlice(queryValues)
		sort.Strings(queryKeys)
		for _, k := range queryKeys {
			if k != "X-Signature" {
				queryPair := fmt.Sprintf("%s=%s", k, queryValues[k][0])
				SignableSlice = append(SignableSlice, queryPair)
			}
		}

		SignableBytes := []byte(strings.Join(SignableSlice, "&"))

		expectedSignature := utils.Sign(key, salt, SignableBytes)
		if expectedSignature != args.Signature {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": "wrong signature",
			})
			c.Abort()
			return
		}

		if args.ExpiresAt != 0 && time.Now().UnixMilli() > args.ExpiresAt {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": "link expired",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
