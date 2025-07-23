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

type CalculateTotalCostController struct {
	useCase usecases.CalculateTotalCostUsecase
	logger  logger.Logger
}

func NewCalculateTotalCostController(
	handler *gin.Engine,
	useCase usecases.CalculateTotalCostUsecase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	ct := &CalculateTotalCostController{
		useCase: useCase,
		logger:  logger,
	}

	handler.POST("/subscriptions/total", ct.CalculateTotalCost, middleware.HandleErrors)
}

// CalculateTotalCost godoc
// @Summary Рассчет общую стоимость подписки
// @Description Расчет общей стоимость подписок за определенный период с использованием дополнительных фильтров
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body requests.CalculateTotalCost true "структура запроса"
// @Success 200 {object} responses.CalculateTotalCost
// @Failure      400 {object} string "некорректный формат запроса"
// @Failure      500 {object} string "внутренняя ошибка сервера"
// @Router /subscriptions/total [post]
func (ct *CalculateTotalCostController) CalculateTotalCost(c *gin.Context) {
	var req requests.CalculateTotalCost
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	response, err := ct.useCase.CalculateTotalCost(c, req)
	if err != nil {
		middleware.AddGinError(c, errors.Wrap(err, "failed to calculate total cost"))

		return
	}
	c.JSON(http.StatusOK, response)
}
