package wallets

import (
	"context"
	"errors"
	"tech_task/internal/repository/wallets"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			u.logger.Error("CreateWalletOperation failed: foreign key violation", zap.Error(err))
			return 0, pgx.ErrNoRows
		}

		u.logger.Error("CreateWalletOperation failed", zap.Error(err))
		return 0, err
	}

	idStr := walletID.String()

	u.wl.Lock(idStr)
	defer u.wl.Unlock(idStr)

	switch operationType {
	case "DEPOSIT":
		err = u.repo.AddAmount(ctx, &wallets.AddAmountParams{
			Amount: walletAmount,
			ID: pgtype.UUID{
				Bytes: walletID,
				Valid: true,
			}})
		if err != nil {
			u.logger.Error("AddAmount failed", zap.Error(err))
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
			u.logger.Error("SubtractAmountIfEnough failed", zap.Error(err))
			return 0, err
		}
	}

	balance, err := u.repo.GetWalletAmount(ctx, &wallets.GetWalletAmountParams{
		ID: pgtype.UUID{
			Bytes: walletID,
			Valid: true,
		},
	})

	u.logger.Info("GetWalletAmount success")
	return balance, nil
}
