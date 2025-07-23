package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	mockSubRepo *MockSubRepository
)

func initDeleteSubTestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSubRepo = NewMockSubRepository(ctrl)
}

func TestDeleteSubscription_Success(t *testing.T) {
	initDeleteSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()

	mockSubRepo.EXPECT().Delete(ctx, subID).Return(nil)

	useCase := NewDeleteSubUsecase(mockSubRepo, mockLogger)
	err := useCase.DeleteSubscription(ctx, subID)

	assert.NoError(t, err)
}

func TestDeleteSubscription_Failure_InvalidSubID(t *testing.T) {
	initDeleteSubTestMocks(t)
	ctx := context.Background()
	subID := "invalid-uuid"

	useCase := NewDeleteSubUsecase(mockSubRepo, mockLogger)
	err := useCase.DeleteSubscription(ctx, subID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidUUID)
}

func TestDeleteSubscription_Failure_NotFound(t *testing.T) {
	initDeleteSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()

	mockSubRepo.EXPECT().Delete(ctx, subID).Return(ErrEntityNotFound)

	useCase := NewDeleteSubUsecase(mockSubRepo, mockLogger)
	err := useCase.DeleteSubscription(ctx, subID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEntityNotFound)
}

func TestDeleteSubscription_Failure_DatabaseError(t *testing.T) {
	initDeleteSubTestMocks(t)
	ctx := context.Background()
	subID := uuid.New().String()

	expectedErr := errors.New("database error")
	mockSubRepo.EXPECT().Delete(ctx, subID).Return(expectedErr)

	useCase := NewDeleteSubUsecase(mockSubRepo, mockLogger)
	err := useCase.DeleteSubscription(ctx, subID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}
