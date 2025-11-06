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
	"github.com/jackc/pgx/v5"
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

func TestUseCase_GetWalletAmount_NoRows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockQuerier(ctrl)
	wl := helper.NewWalletLocker()

	u := wallets_usecase.NewUseCase(zaptest.NewLogger(t), mockRepo, wl)

	walletID := uuid.New()

	mockRepo.EXPECT().
		GetWalletAmount(gomock.Any(), gomock.Any()).
		Return(float32(0), errors.New("some unknown error")).
		Times(1)

	_, err := u.GetWalletAmount(context.Background(), walletID)
	if !errors.Is(err, errors.New("some unknown error")) && !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("expected pgx.ErrNoRows or original error, got %v", err)
	}
}

func TestUseCase_GetWalletAmount_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockQuerier(ctrl)
	wl := helper.NewWalletLocker()

	u := wallets_usecase.NewUseCase(zaptest.NewLogger(t), mockRepo, wl)

	walletID := uuid.New()
	repoErr := errors.New("no rows in result set")

	mockRepo.EXPECT().
		GetWalletAmount(gomock.Any(), gomock.Any()).
		Return(float32(0), repoErr).
		Times(1)

	_, err := u.GetWalletAmount(context.Background(), walletID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != repoErr.Error() {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
