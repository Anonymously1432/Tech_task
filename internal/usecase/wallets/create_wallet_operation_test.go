package wallets_test

import (
	"context"
	"errors"
	"tech_task/internal/helper"
	"tech_task/internal/repository/wallets"
	mockrepo "tech_task/internal/repository/wallets/mocks"
	wallets_usecase "tech_task/internal/usecase/wallets"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap/zaptest"
)

func TestUseCase_CreateWalletOperation_DepositSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockQuerier(ctrl)
	wl := helper.NewWalletLocker()
	u := wallets_usecase.NewUseCase(zaptest.NewLogger(t), mockRepo, wl)

	ctx := context.Background()
	walletID := uuid.New()
	amount := float32(100)

	mockRepo.EXPECT().
		CreateWalletOperation(ctx, &wallets.CreateWalletOperationParams{
			Operationtype: "DEPOSIT",
			WalletID:      pgtype.UUID{Bytes: walletID, Valid: true},
			Amount:        amount,
		}).
		Return(&wallets.RequestsHistory{}, nil)

	mockRepo.EXPECT().
		AddAmount(ctx, &wallets.AddAmountParams{
			Amount: amount,
			ID:     pgtype.UUID{Bytes: walletID, Valid: true},
		}).
		Return(nil)

	expectedBalance := float32(500)
	mockRepo.EXPECT().
		GetWalletAmount(ctx, &wallets.GetWalletAmountParams{
			ID: pgtype.UUID{Bytes: walletID, Valid: true},
		}).
		Return(expectedBalance, nil)

	balance, err := u.CreateWalletOperation(ctx, walletID, "DEPOSIT", amount)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if balance != expectedBalance {
		t.Fatalf("expected %v, got %v", expectedBalance, balance)
	}
}

func TestUseCase_CreateWalletOperation_WithdrawSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockQuerier(ctrl)
	wl := helper.NewWalletLocker()
	u := wallets_usecase.NewUseCase(zaptest.NewLogger(t), mockRepo, wl)

	ctx := context.Background()
	walletID := uuid.New()
	amount := float32(50)

	mockRepo.EXPECT().
		CreateWalletOperation(ctx, &wallets.CreateWalletOperationParams{
			Operationtype: "WITHDRAW",
			WalletID:      pgtype.UUID{Bytes: walletID, Valid: true},
			Amount:        amount,
		}).
		Return(&wallets.RequestsHistory{}, nil)

	mockRepo.EXPECT().
		SubtractAmountIfEnough(ctx, &wallets.SubtractAmountIfEnoughParams{
			Amount: amount,
			ID:     pgtype.UUID{Bytes: walletID, Valid: true},
		}).
		Return(nil)

	expectedBalance := float32(450)
	mockRepo.EXPECT().
		GetWalletAmount(ctx, &wallets.GetWalletAmountParams{
			ID: pgtype.UUID{Bytes: walletID, Valid: true},
		}).
		Return(expectedBalance, nil)

	balance, err := u.CreateWalletOperation(ctx, walletID, "WITHDRAW", amount)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if balance != expectedBalance {
		t.Fatalf("expected %v, got %v", expectedBalance, balance)
	}
}

func TestUseCase_CreateWalletOperation_ErrorCreateOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockQuerier(ctrl)
	wl := helper.NewWalletLocker()
	u := wallets_usecase.NewUseCase(zaptest.NewLogger(t), mockRepo, wl)

	ctx := context.Background()
	walletID := uuid.New()

	mockRepo.EXPECT().
		CreateWalletOperation(gomock.Any(), gomock.Any()).
		Return(&wallets.RequestsHistory{}, errors.New("create error"))

	_, err := u.CreateWalletOperation(ctx, walletID, "DEPOSIT", 100)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUseCase_CreateWalletOperation_ErrorAddAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockQuerier(ctrl)
	wl := helper.NewWalletLocker()
	u := wallets_usecase.NewUseCase(zaptest.NewLogger(t), mockRepo, wl)

	ctx := context.Background()
	walletID := uuid.New()

	amount := float32(100)

	mockRepo.EXPECT().
		CreateWalletOperation(gomock.Any(), gomock.Any()).
		Return(&wallets.RequestsHistory{}, nil)

	mockRepo.EXPECT().
		AddAmount(gomock.Any(), gomock.Any()).
		Return(errors.New("add error"))

	_, err := u.CreateWalletOperation(ctx, walletID, "DEPOSIT", amount)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUseCase_CreateWalletOperation_ErrorSubtractAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockQuerier(ctrl)
	wl := helper.NewWalletLocker()
	u := wallets_usecase.NewUseCase(zaptest.NewLogger(t), mockRepo, wl)

	ctx := context.Background()
	walletID := uuid.New()
	amount := float32(50)

	mockRepo.EXPECT().
		CreateWalletOperation(gomock.Any(), gomock.Any()).
		Return(&wallets.RequestsHistory{}, nil)

	mockRepo.EXPECT().
		SubtractAmountIfEnough(gomock.Any(), gomock.Any()).
		Return(errors.New("sub error"))

	_, err := u.CreateWalletOperation(ctx, walletID, "WITHDRAW", amount)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUseCase_CreateWalletOperation_ForeignKeyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockQuerier(ctrl)
	wl := helper.NewWalletLocker()
	u := wallets_usecase.NewUseCase(zaptest.NewLogger(t), mockRepo, wl)

	ctx := context.Background()
	walletID := uuid.New()

	pgErr := &pgconn.PgError{Code: "23503"}

	mockRepo.EXPECT().
		CreateWalletOperation(ctx, gomock.Any()).
		Return(nil, pgErr)

	_, err := u.CreateWalletOperation(ctx, walletID, "DEPOSIT", 100)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}
