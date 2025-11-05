package wallets

import (
	"context"
	"tech_task/internal/repository/wallets"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IUseCase interface {
	GetWalletAmount(ctx context.Context, walletID uuid.UUID) (float32, error)
	CreateWalletOperation() error
}

type UseCase struct {
	logger *zap.Logger
	repo   *wallets.Queries
}

func NewUseCase(logger *zap.Logger, repo *wallets.Queries) IUseCase {
	return &UseCase{
		logger: logger,
		repo:   repo,
	}
}
