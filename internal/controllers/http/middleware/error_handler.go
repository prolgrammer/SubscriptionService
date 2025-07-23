package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"subscription_service/internal/controllers"
	"subscription_service/internal/usecases"
)

func (m *middleware) HandleErrors(c *gin.Context) {
	if len(c.Errors) > 0 {
		err := c.Errors.Last()

		if errors.Is(err, controllers.ErrDataBindError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, usecases.ErrEntityAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, usecases.ErrEntityNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
			return
		}

		m.logger.Err(err).Error().Msgf("Unexpected error: ")
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Internal server error")
	}
}
