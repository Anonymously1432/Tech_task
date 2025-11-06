package domain

import (
	"errors"

	"github.com/google/uuid"
)

type WalletOperationRequest struct {
	WalletId      uuid.UUID `json:"valletId"`
	OperationType string    `json:"operationType"`
	Amount        float32   `json:"amount"`
}

func (r *WalletOperationRequest) Validate() error {
	if r.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if r.OperationType != "DEPOSIT" && r.OperationType != "WITHDRAW" {
		return errors.New("operationType must be 'deposit' or 'withdraw'")
	}

	return nil
}

type WalletOperationResponse struct {
	WalletId uuid.UUID `json:"valletId"`
	Amount   float32   `json:"amount"`
}

type GetWalletAmountResponse float32
