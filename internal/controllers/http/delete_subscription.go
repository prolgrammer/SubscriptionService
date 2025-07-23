package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"subscription_service/internal/controllers"
	"subscription_service/internal/controllers/http/middleware"
	"subscription_service/internal/usecases"
	"subscription_service/pkg/logger"
)

type deleteSubController struct {
	useCase usecases.DeleteSubUseCase
	logger  logger.Logger
}

func NewDeleteSubController(
	handler *gin.Engine,
	useCase usecases.DeleteSubUseCase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	ct := &deleteSubController{
		useCase: useCase,
		logger:  logger,
	}

	handler.DELETE("/subscriptions/:sub_id", ct.DeleteSubscription, middleware.HandleErrors)
}

// DeleteSubscription godoc
// @Summary Удаление подписки
// @Description Удаление подписки по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "path format"
// @Success 200
// @Failure 	 400 {object} string "некорректный формат запроса"
// @Failure      404 {object} string "подписка не найдена"
// @Failure      500 {object} string "внутренняя ошибка сервера
// @Router /subscriptions/{sub_id} [delete]
func (ds *deleteSubController) DeleteSubscription(c *gin.Context) {
	subId := c.Param("sub_id")
	if subId == "" {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	if err := ds.useCase.DeleteSubscription(c, subId); err != nil {
		middleware.AddGinError(c, errors.Wrap(err, "failed to delete subscription"))
		return
	}

	c.Status(http.StatusOK)
}
