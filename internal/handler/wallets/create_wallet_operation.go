package wallets

import (
	"tech_task/internal/domain"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *Handler) CreateWalletOperation(c *fiber.Ctx) error {
	req := new(domain.WalletOperationRequest)

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("CreateWalletOperation: body parsing error", zap.Any("err", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := req.Validate(); err != nil {
		h.logger.Error("CreateWalletOperation: validation error", zap.Any("err", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	newAmount, err := h.Uc.CreateWalletOperation(c.Context(), req.WalletId, req.OperationType, req.Amount)
	if err != nil {
		h.logger.Error("Error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info("Create wallet operation success")
	return c.Status(fiber.StatusOK).JSON(domain.WalletOperationResponse{
		WalletId: req.WalletId,
		Amount:   newAmount,
	})
}
