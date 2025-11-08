package middleware

import (
	"net/http"

	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}


func ValidateRequest(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				util.NewErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		if err := validate.Struct(obj); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				util.NewErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		c.Next()
	}
}


func ValidateQueryParams(requiredParams ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, param := range requiredParams {
			if c.Query(param) == "" {
				c.AbortWithStatusJSON(http.StatusBadRequest,
					util.NewErrorResponse(http.StatusBadRequest, "Required query parameter '"+param+"' is missing"))
				return
			}
		}
		c.Next()
	}
}
