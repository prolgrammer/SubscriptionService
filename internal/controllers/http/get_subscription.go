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

type getSubController struct {
	useCase usecases.GetSubUsecase
	logger  logger.Logger
}

func NewGetSubController(
	handler *gin.Engine,
	useCase usecases.GetSubUsecase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	ct := &getSubController{
		useCase: useCase,
		logger:  logger,
	}

	handler.GET("/subscriptions/:sub_id", ct.GetSubscription, middleware.HandleErrors)
}

// GetSubscription godoc
// @Summary Запрос на получение подписки
// @Description Запрос на получение подписки по ее ID
// @Tags subscriptions
// @Produce      json
// @Param 	     id path string true "path format"
// @Success 	 200 {object} responses.SubResponse
// @Failure 	 400 {object} string "некорректный формат запроса"
// @Failure      404 {object} string "подписка не найдена"
// @Failure      500 {object} string "внутренняя ошибка сервера
// @Router /subscriptions/{sub_id} [get]
func (gs *getSubController) GetSubscription(c *gin.Context) {
	subId := c.Param("sub_id")
	if subId == "" {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	response, err := gs.useCase.GetSubscription(c, subId)
	if err != nil {
		middleware.AddGinError(c, errors.Wrap(err, "failed to get subscription"))
		return
	}

	c.JSON(http.StatusOK, response)
}
