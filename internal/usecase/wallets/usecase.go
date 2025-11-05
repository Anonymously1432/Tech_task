package wallets

import (
	"go.uber.org/zap"
)

type IUseCase interface {
	Test() error
}

type UseCase struct {
	logger *zap.Logger
	//repo   *wallets.Queries
}

func NewUseCase(logger *zap.Logger /* repo *wallets.Queries */) IUseCase {
	return &UseCase{
		logger: logger,
		//repo:   repo,
	}
}
