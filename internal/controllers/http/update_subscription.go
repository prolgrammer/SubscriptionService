package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"subscription_service/internal/controllers"
	"subscription_service/internal/controllers/http/middleware"
	"subscription_service/internal/controllers/requests"
	"subscription_service/internal/usecases"
	"subscription_service/pkg/logger"
)

type updateSubController struct {
	useCase usecases.UpdateSubUseCase
	logger  logger.Logger
}

func NewUpdateSubController(
	handler *gin.Engine,
	useCase usecases.UpdateSubUseCase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	ct := &updateSubController{
		useCase: useCase,
		logger:  logger,
	}

	handler.PUT("/subscriptions/:sub_id", ct.UpdateSubscription, middleware.HandleErrors)
}

// UpdateSubscription godoc
// @Summary Обновление подписки
// @Description Обновление подписки по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param sub_id path string true "path format"
// @Param subscription body requests.SubRequest true "структура запроса"
// @Success 	 200 {object} responses.SubResponse
// @Failure 	 400 {object} string "некорректный формат запроса"
// @Failure      404 {object} string "подписка не найдена"
// @Failure      500 {object} string "внутренняя ошибка сервера
// @Router /subscriptions/{sub_id} [put]
func (us *updateSubController) UpdateSubscription(c *gin.Context) {
	subId := c.Param("sub_id")
	if subId == "" {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	var req requests.SubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	response, err := us.useCase.UpdateSubscription(c, subId, req)
	if err != nil {
		middleware.AddGinError(c, errors.Wrap(err, "failed to update subscription"))
		return
	}

	c.JSON(http.StatusOK, response)
}
