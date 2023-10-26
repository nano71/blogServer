package api

import (
	"blogServer/database"
	"blogServer/response"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

func isParameterMissingError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "failed on the 'required'")
}
func preprocess(c *gin.Context, p interface{}, callback func(*gorm.DB)) {
	c.Header("Content-Type", "application/json")

	if err := c.ShouldBind(p); err != nil {
		if isParameterMissingError(err) {
			response.MissingParameters(c)
		} else {
			response.ParameterError(c)
		}
		return
	}

	callback(database.GetDB())

}

func GetPermission(c *gin.Context) {
	//p := &struct {
	//	SecretKey string `json:"secretKey" binding:"required"`
	//}{}
}

func TimeoutMiddleware(limit time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(limit*time.Millisecond),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(response.Timeout),
	)
}
