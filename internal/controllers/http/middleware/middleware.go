package middleware

import (
	"github.com/gin-gonic/gin"
	"subscription_service/pkg/logger"
)

type middleware struct {
	logger logger.Logger
}

type Middleware interface {
	HandleErrors(c *gin.Context)
}

func NewMiddleware(
	logger logger.Logger) Middleware {
	return &middleware{logger: logger}
}
