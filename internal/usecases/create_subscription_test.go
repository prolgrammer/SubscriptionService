package usecases

import (
	"context"
	"github.com/pkg/errors"
	"subscription_service/internal/controllers/requests"
	"subscription_service/internal/entities"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	mockSubRepo *MockSubRepository
)

func initCreateSubTestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSubRepo = NewMockSubRepository(ctrl)
}

func TestCreateSubscription_Success(t *testing.T) {
	initCreateSubTestMocks(t)
	ctx := context.Background()
	req := requests.SubRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
	}

	sub := gomock.AssignableToTypeOf(&entities.Subscription{})
	mockSubRepo.EXPECT().Insert(ctx, sub).Return(nil)

	useCase := NewCreateSubUsecase(mockSubRepo, mockLogger)
	response, err := useCase.CreateSubscription(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.ServiceName, response.ServiceName)
	assert.Equal(t, req.Price, response.Price)
	assert.Equal(t, req.UserID, response.UserID)
	assert.Equal(t, req.StartDate, response.StartDate)
}

func TestCreateSubscription_Failure_InvalidUserID(t *testing.T) {
	initCreateSubTestMocks(t)
	ctx := context.Background()
	req := requests.SubRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "invalid-uuid",
		StartDate:   "07-2025",
	}

	useCase := NewCreateSubUsecase(mockSubRepo, mockLogger)
	_, err := useCase.CreateSubscription(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidUUID)
}

func TestCreateSubscription_Failure_InvalidStartDate(t *testing.T) {
	initCreateSubTestMocks(t)
	ctx := context.Background()
	req := requests.SubRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New().String(),
		StartDate:   "invalid-date",
	}

	useCase := NewCreateSubUsecase(mockSubRepo)
	_, err := useCase.CreateSubscription(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidDateFormat)
}

func TestCreateSubscription_Failure_InsertError(t *testing.T) {
	initCreateSubTestMocks(t)
	ctx := context.Background()
	req := requests.SubRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
	}

	expectedErr := errors.New("database error")
	mockSubRepo.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

	useCase := NewCreateSubUsecase(mockSubRepo)
	_, err := useCase.CreateSubscription(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}
