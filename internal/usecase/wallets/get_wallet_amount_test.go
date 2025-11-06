package wallets_test

import (
	"context"
	"errors"
	"tech_task/internal/helper"
	"tech_task/internal/repository/wallets"
	wallets_usecase "tech_task/internal/usecase/wallets"

	"testing"

	mockrepo "tech_task/internal/repository/wallets/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap/zaptest"
)

func TestUseCase_GetWalletAmount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockQuerier(ctrl)
	wl := helper.NewWalletLocker()

	u := wallets_usecase.NewUseCase(zaptest.NewLogger(t), mockRepo, wl)

	walletID := uuid.New()
	expectedAmount := float32(150.75)

	mockRepo.EXPECT().
		GetWalletAmount(
			gomock.Any(),
			&wallets.GetWalletAmountParams{
				ID: pgtype.UUID{Bytes: walletID, Valid: true},
			},
		).
		Return(expectedAmount, nil).
		Times(1)

	amount, err := u.GetWalletAmount(context.Background(), walletID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if amount != expectedAmount {
		t.Fatalf("expected %v, got %v", expectedAmount, amount)
	}
}

func TestUseCase_GetWalletAmount_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockQuerier(ctrl)
	wl := helper.NewWalletLocker()

	u := wallets_usecase.NewUseCase(zaptest.NewLogger(t), mockRepo, wl)

	walletID := uuid.New()
	expectedErr := errors.New("repository error")

	mockRepo.EXPECT().
		GetWalletAmount(
			gomock.Any(),
			gomock.Any(),
		).
		Return(float32(0), expectedErr).
		Times(1)

	_, err := u.GetWalletAmount(context.Background(), walletID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
