package wallets

import (
	"tech_task/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *Handler) GetWalletAmount(c *fiber.Ctx) error {
	walletID := c.Params("id")
	walletUUID, err := uuid.Parse(walletID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	amount, err := h.Uc.GetWalletAmount(c.Context(), walletUUID)
	if err != nil {
		h.logger.Error("Error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info("Success")
	return c.Status(fiber.StatusOK).JSON(domain.GetWalletAmountResponse(amount))
}
