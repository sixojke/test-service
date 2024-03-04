package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func newErrorResponse(c *gin.Context, statusCode int, r errorResponse) {
	logrus.Error(r)
	c.AbortWithStatusJSON(statusCode, r)
}
