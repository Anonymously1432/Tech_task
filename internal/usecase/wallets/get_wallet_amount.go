package wallets

import (
	"context"
	"tech_task/internal/repository/wallets"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (u *UseCase) GetWalletAmount(ctx context.Context, walletID uuid.UUID) (float32, error) {
	amount, err := u.repo.GetWalletAmount(ctx, &wallets.GetWalletAmountParams{ID: pgtype.UUID{
		Bytes: walletID,
		Valid: true,
	}})
	if err != nil {
		return 0, err
	}

	return amount, nil
}
