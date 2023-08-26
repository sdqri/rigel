package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdqri/rigel/utils"
)

type XParams struct {
	Signature string `form:"X-Signature" binding:"required"`
	ExpiresAt int64  `form:"X-ExpiresAt"`
}

func NewSignatureValidator(prefix string, signatory *utils.Signatory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var args XParams
		err := c.BindQuery(&args)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		expectedSignature := signatory.SignURL(*c.Request.URL)
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
