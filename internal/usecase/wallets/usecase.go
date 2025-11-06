package wallets

import (
	"context"
	"tech_task/internal/helper"
	"tech_task/internal/repository/wallets"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IUseCase interface {
	GetWalletAmount(ctx context.Context, walletID uuid.UUID) (float32, error)
	CreateWalletOperation(ctx context.Context, walletID uuid.UUID, operationType string, walletAmount float32) (float32, error)
}

type UseCase struct {
	logger *zap.Logger
	repo   *wallets.Queries
	wl     *helper.WalletLocker
}

func NewUseCase(logger *zap.Logger, repo *wallets.Queries, wl *helper.WalletLocker) IUseCase {
	return &UseCase{
		logger: logger,
		repo:   repo,
		wl:     wl,
	}
}
