package wallets

import (
	"tech_task/internal/usecase/wallets"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	Uc     wallets.IUseCase
	logger *zap.Logger
}

func NewHandler(log *zap.Logger, uc wallets.IUseCase) *Handler {
	return &Handler{Uc: uc, logger: log}
}

func RegisterRoutes(app fiber.Router, h *Handler) {
	app.Get("/", h.Test)
}
