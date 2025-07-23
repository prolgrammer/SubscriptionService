package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"subscription_service/internal/entities"
	"subscription_service/pkg/logger"
)

var (
	mockSubRepo *MockSubRepository
)

func initGetListSubTestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSubRepo = NewMockSubRepository(ctrl)
}

func TestGetListSubscriptions_Success(t *testing.T) {
	initGetListSubTestMocks(t)
	ctx := context.Background()
	limit, offset := 10, 0
	startDate, _ := time.Parse("2006-01-02", "2025-07-01")
	endDate, _ := time.Parse("2006-01-02", "2025-12-01")

	mockSubs := []entities.Subscription{
		{
			ID:          uuid.New(),
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      uuid.New(),
			StartDate:   startDate,
			EndDate:     &endDate,
		},
		{
			ID:          uuid.New(),
			ServiceName: "Spotify",
			Price:       300,
			UserID:      uuid.New(),
			StartDate:   startDate,
			EndDate:     nil,
		},
	}

	mockSubRepo.EXPECT().SelectAll(ctx, limit, offset).Return(mockSubs, nil)

	useCase := NewGetListSubUsecase(mockSubRepo, mockLogger)
	response, err := useCase.GetListSubscriptions(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response, 2)
	assert.Equal(t, mockSubs[0].ServiceName, response[0].ServiceName)
	assert.Equal(t, mockSubs[0].Price, response[0].Price)
	assert.Equal(t, mockSubs[0].UserID.String(), response[0].UserID)
	assert.Equal(t, "07-2025", response[0].StartDate)
	assert.Equal(t, "12-2025", response[0].EndDate)
	assert.Equal(t, mockSubs[1].ServiceName, response[1].ServiceName)
	assert.Equal(t, mockSubs[1].Price, response[1].Price)
	assert.Equal(t, mockSubs[1].UserID.String(), response[1].UserID)
	assert.Equal(t, "07-2025", response[1].StartDate)
	assert.Empty(t, response[1].EndDate)
}

func TestGetListSubscriptions_Failure_DatabaseError(t *testing.T) {
	initGetListSubTestMocks(t)
	ctx := context.Background()
	limit, offset := 10, 0

	expectedErr := errors.New("database error")
	mockSubRepo.EXPECT().SelectAll(ctx, limit, offset).Return(nil, expectedErr)

	useCase := NewGetListSubUsecase(mockSubRepo, mockLogger)
	_, err := useCase.GetListSubscriptions(ctx, limit, offset)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}
