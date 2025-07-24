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
)

var (
	mockGetSubRepo *MockGetSubRepository
)

func initGetSubTestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockGetSubRepo = NewMockGetSubRepository(ctrl)
}

func TestGetSubscription_Success(t *testing.T) {
	initGetSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()
	startDate, _ := time.Parse("2006-01-02", "2025-07-01")
	endDate, _ := time.Parse("2006-01-02", "2025-12-01")

	mockSub := entities.Subscription{
		ID:          uuid.MustParse(subID),
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   startDate,
		EndDate:     &endDate,
	}

	mockGetSubRepo.EXPECT().SelectByID(ctx, subID).Return(mockSub, nil)

	useCase := NewGetSubUseCase(mockGetSubRepo, mockLogger)
	response, err := useCase.GetSubscription(ctx, subID)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, mockSub.ID.String(), response.ID)
	assert.Equal(t, mockSub.ServiceName, response.ServiceName)
	assert.Equal(t, mockSub.Price, response.Price)
	assert.Equal(t, mockSub.UserID.String(), response.UserID)
	assert.Equal(t, "07-2025", response.StartDate)
	assert.Equal(t, "12-2025", response.EndDate)
}

func TestGetSubscription_Failure_InvalidSubID(t *testing.T) {
	initGetSubTestMocks(t)
	ctx := context.Background()
	subID := "invalid-uuid"

	useCase := NewGetSubUseCase(mockGetSubRepo, mockLogger)
	_, err := useCase.GetSubscription(ctx, subID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidUUID)
}

func TestGetSubscription_Failure_NotFound(t *testing.T) {
	initGetSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()

	mockGetSubRepo.EXPECT().SelectByID(ctx, subID).Return(entities.Subscription{}, ErrEntityNotFound)

	useCase := NewGetSubUseCase(mockGetSubRepo, mockLogger)
	_, err := useCase.GetSubscription(ctx, subID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEntityNotFound)
}

func TestGetSubscription_Failure_DatabaseError(t *testing.T) {
	initGetSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()

	expectedErr := errors.New("database error")
	mockGetSubRepo.EXPECT().SelectByID(ctx, subID).Return(entities.Subscription{}, expectedErr)

	useCase := NewGetSubUseCase(mockGetSubRepo, mockLogger)
	_, err := useCase.GetSubscription(ctx, subID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}
