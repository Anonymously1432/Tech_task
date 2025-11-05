package domain

type WalletOperationRequest struct {
	WalletId      string  `json:"valletId"`
	OperationType string  `json:"operationType"`
	Amount        float32 `json:"amount"`
}

type WalletOperationResponse struct{}

type GetWalletAmountResponse float32
