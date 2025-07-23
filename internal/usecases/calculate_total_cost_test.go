package usecases

import (
	"context"
	"errors"
	"net/http"
	"subscription_service/pkg/logger"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"subscription_service/internal/controllers/requests"
)

var (
	mockLogger           *logger.MockLogger
	mockCalculateSubRepo *MockCalculateTotalCostRepository
)

func initCalculateTotalCostTestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCalculateSubRepo = NewMockCalculateTotalCostRepository(ctrl)
	mockLogger = logger.NewMockLogger(ctrl)
}

func TestCalculateTotalCost_Success(t *testing.T) {
	initCalculateTotalCostTestMocks(t)
	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	startPeriod, _ := time.Parse("2006-01-02", "2025-07-01")
	endPeriod, _ := time.Parse("2006-01-02", "2025-12-01")
	userID := uuid.New().String()
	serviceName := "Yandex Plus"

	req := requests.CalculateTotalCost{
		StartPeriod: "07-2025",
		EndPeriod:   "12-2025",
		UserID:      userID,
		ServiceName: serviceName,
	}

	mockCalculateSubRepo.EXPECT().CalculateTotalCost(ctx, startPeriod, endPeriod, &userID, &serviceName).Return(400, nil)

	useCase := NewCalculateTotalCostUseCase(mockCalculateSubRepo, mockLogger)
	response, err := useCase.CalculateTotalCost(c, req)

	assert.NoError(t, err)
	assert.Equal(t, 400, response.Total)
}

func TestCalculateTotalCost_Success_NoFilters(t *testing.T) {
	initCalculateTotalCostTestMocks(t)
	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	startPeriod, _ := time.Parse("2006-01-02", "2025-07-01")
	endPeriod, _ := time.Parse("2006-01-02", "2025-12-01")

	req := requests.CalculateTotalCost{
		StartPeriod: "07-2025",
		EndPeriod:   "12-2025",
	}

	mockCalculateSubRepo.EXPECT().CalculateTotalCost(ctx, startPeriod, endPeriod, nil, nil).Return(700, nil)

	useCase := NewCalculateTotalCostUseCase(mockCalculateSubRepo, mockLogger)
	response, err := useCase.CalculateTotalCost(c, req)

	assert.NoError(t, err)
	assert.Equal(t, 700, response.Total)
}

func TestCalculateTotalCost_Failure_InvalidStartPeriod(t *testing.T) {
	initCalculateTotalCostTestMocks(t)
	c, _ := gin.CreateTestContext(nil)
	req := requests.CalculateTotalCost{
		StartPeriod: "invalid-date",
		EndPeriod:   "12-2025",
	}

	useCase := NewCalculateTotalCostUseCase(mockCalculateSubRepo, mockLogger)
	_, err := useCase.CalculateTotalCost(c, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidDateFormat)
}

func TestCalculateTotalCost_Failure_InvalidEndPeriod(t *testing.T) {
	initCalculateTotalCostTestMocks(t)
	c, _ := gin.CreateTestContext(nil)
	req := requests.CalculateTotalCost{
		StartPeriod: "07-2025",
		EndPeriod:   "invalid-date",
	}

	useCase := NewCalculateTotalCostUseCase(mockCalculateSubRepo, mockLogger)
	_, err := useCase.CalculateTotalCost(c, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidDateFormat)
}

func TestCalculateTotalCost_Failure_InvalidUserID(t *testing.T) {
	initCalculateTotalCostTestMocks(t)
	c, _ := gin.CreateTestContext(nil)
	req := requests.CalculateTotalCost{
		StartPeriod: "07-2025",
		EndPeriod:   "12-2025",
		UserID:      "invalid-uuid",
	}

	useCase := NewCalculateTotalCostUseCase(mockCalculateSubRepo, mockLogger)
	_, err := useCase.CalculateTotalCost(c, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidUUID)
}

func TestCalculateTotalCost_Failure_DatabaseError(t *testing.T) {
	initCalculateTotalCostTestMocks(t)
	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	startPeriod, _ := time.Parse("2006-01-02", "2025-07-01")
	endPeriod, _ := time.Parse("2006-01-02", "2025-12-01")

	req := requests.CalculateTotalCost{
		StartPeriod: "07-2025",
		EndPeriod:   "12-2025",
	}

	expectedErr := errors.New("database error")
	mockCalculateSubRepo.EXPECT().CalculateTotalCost(ctx, startPeriod, endPeriod, nil, nil).Return(700, nil)

	useCase := NewCalculateTotalCostUseCase(mockCalculateSubRepo, mockLogger)
	_, err := useCase.CalculateTotalCost(c, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}
