package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"subscription_service/internal/controllers"
	"subscription_service/internal/controllers/http/middleware"
	"subscription_service/internal/usecases"
	"subscription_service/pkg/logger"
)

type getListSubController struct {
	useCase usecases.GetListSubUsecase
	logger  logger.Logger
}

func NewGetListSubController(
	handler *gin.Engine,
	useCase usecases.GetListSubUsecase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	ct := &getListSubController{
		useCase: useCase,
		logger:  logger,
	}

	handler.GET("/subscriptions", ct.GetListSubscriptions, middleware.HandleErrors)
}

// GetListSubscriptions godoc
// @Summary Получение списка подписок
// @Description Возвращает список подписок с поддержкой пагинации
// @Produce      json
// @Param limit query int false "Количество подписок на странице" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success      200 {object} []responses.SubResponse
// @Failure      400 {object} string "некорректный формат запроса"
// @Failure      500 {object} string "внутренняя ошибка сервера"
// @Router /subscriptions [get]
func (gl *getListSubController) GetListSubscriptions(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		middleware.AddGinError(c, controllers.ErrInvalidPaginationParams)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		middleware.AddGinError(c, controllers.ErrInvalidPaginationParams)
		return
	}

	response, err := gl.useCase.GetListSubscriptions(limit, offset)
	if err != nil {
		middleware.AddGinError(c, errors.Wrap(err, "failed to get list subscription"))

		return
	}
	c.JSON(http.StatusOK, response)
}
