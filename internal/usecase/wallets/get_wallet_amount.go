package wallets

import (
	"context"
	"tech_task/internal/repository/wallets"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (u *UseCase) GetWalletAmount(ctx context.Context, walletID uuid.UUID) (float32, error) {
	amount, err := u.repo.GetWalletAmount(ctx, &wallets.GetWalletAmountParams{ID: pgtype.UUID{
		Bytes: walletID,
		Valid: true,
	}})
	if err != nil {
		u.logger.Error("GetWalletAmount error", zap.Error(err))
		return 0, err
	}

	u.logger.Info("GetWalletAmount success", zap.Any("amount", amount))
	return amount, nil
}
