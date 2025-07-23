package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"subscription_service/internal/controllers/requests"
	"subscription_service/internal/entities"
)

var (
	mockSubRepo *MockSubRepository
)

func initUpdateSubTestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSubRepo = NewMockSubRepository(ctrl)
}

func TestUpdateSubscription_Success(t *testing.T) {
	initUpdateSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()
	req := requests.SubRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
	}

	sub := gomock.AssignableToTypeOf(&entities.Subscription{})
	mockSubRepo.EXPECT().Update(ctx, sub).Return(nil)

	useCase := NewUpdateSubUsecase(mockSubRepo, mockLogger)
	response, err := useCase.UpdateSubscription(ctx, subID, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.ServiceName, response.ServiceName)
	assert.Equal(t, req.Price, response.Price)
	assert.Equal(t, req.UserID, response.UserID)
	assert.Equal(t, req.StartDate, response.StartDate)
}

func TestUpdateSubscription_Failure_InvalidSubID(t *testing.T) {
	initUpdateSubTestMocks(t)
	ctx := context.Background()
	req := requests.SubRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
	}

	useCase := NewUpdateSubUsecase(mockSubRepo, mockLogger)
	_, err := useCase.UpdateSubscription(ctx, "invalid-uuid", req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidUUID)
}

func TestUpdateSubscription_Failure_InvalidUserID(t *testing.T) {
	initUpdateSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()
	req := requests.SubRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "invalid-uuid",
		StartDate:   "07-2025",
	}

	useCase := NewUpdateSubUsecase(mockSubRepo, mockLogger)
	_, err := useCase.UpdateSubscription(ctx, subID, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidUUID)
}

func TestUpdateSubscription_Failure_InvalidStartDate(t *testing.T) {
	initUpdateSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()
	req := requests.SubRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New().String(),
		StartDate:   "invalid-date",
	}

	useCase := NewUpdateSubUsecase(mockSubRepo, mockLogger)
	_, err := useCase.UpdateSubscription(ctx, subID, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidDateFormat)
}

func TestUpdateSubscription_Failure_NotFound(t *testing.T) {
	initUpdateSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()
	req := requests.SubRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
	}

	mockSubRepo.EXPECT().Update(ctx, gomock.Any()).Return(ErrEntityNotFound)

	useCase := NewUpdateSubUsecase(mockSubRepo, mockLogger)
	_, err := useCase.UpdateSubscription(ctx, subID, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEntityNotFound)
}

func TestUpdateSubscription_Failure_DatabaseError(t *testing.T) {
	initUpdateSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()
	req := requests.SubRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
	}

	expectedErr := errors.New("database error")
	mockSubRepo.EXPECT().Update(ctx, gomock.Any()).Return(expectedErr)

	useCase := NewUpdateSubUsecase(mockSubRepo, mockLogger)
	_, err := useCase.UpdateSubscription(ctx, subID, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}
