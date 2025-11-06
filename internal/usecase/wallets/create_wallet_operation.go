package wallets

import (
	"context"
	"tech_task/internal/repository/wallets"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (u *UseCase) CreateWalletOperation(ctx context.Context, walletID uuid.UUID, operationType string, walletAmount float32) (float32, error) {
	_, err := u.repo.CreateWalletOperation(ctx, &wallets.CreateWalletOperationParams{
		Operationtype: operationType,
		WalletID: pgtype.UUID{
			Bytes: walletID,
			Valid: true,
		},
		Amount: walletAmount,
	})
	if err != nil {
		return 0, err
	}

	switch operationType {
	case "DEPOSIT":
		err = u.repo.AddAmount(ctx, &wallets.AddAmountParams{
			Amount: walletAmount,
			ID: pgtype.UUID{
				Bytes: walletID,
				Valid: true,
			}})
		if err != nil {
			return 0, err
		}

	case "WITHDRAW":
		err = u.repo.SubtractAmountIfEnough(ctx, &wallets.SubtractAmountIfEnoughParams{
			Amount: walletAmount,
			ID: pgtype.UUID{
				Bytes: walletID,
				Valid: true,
			},
		})
		if err != nil {
			return 0, err
		}
	}

	balance, err := u.repo.GetWalletAmount(ctx, &wallets.GetWalletAmountParams{
		ID: pgtype.UUID{
			Bytes: walletID,
			Valid: true,
		},
	})

	return balance, nil
}
