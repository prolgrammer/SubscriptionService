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

type createSubController struct {
	useCase usecases.CreateSubUsecase
	logger  logger.Logger
}

func NewCreateSubController(
	handler *gin.Engine,
	useCase usecases.CreateSubUsecase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	ct := &createSubController{
		useCase: useCase,
		logger:  logger,
	}

	handler.POST("/subscriptions", ct.CreateSubscription, middleware.HandleErrors)
}

// CreateSubscription godoc
// @Summary Создание подписки
// @Description Создание новой записи подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body requests.SubRequest true "структура запроса"
// @Success 201 {object} responses.SubResponse
// @Failure 400 {object} string "некорректный формат запроса"
// @Failure 500 {object} string "внутренняя ошибка сервера"
// @Router /subscriptions [post]
func (cs *createSubController) CreateSubscription(c *gin.Context) {
	var req requests.SubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	response, err := cs.useCase.CreateSubscription(req)
	if err != nil {
		middleware.AddGinError(c, errors.Wrap(err, "failed to create subscription"))
		return
	}

	c.JSON(http.StatusCreated, response)
}
