package wallets

import (
	"errors"
	"tech_task/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (h *Handler) GetWalletAmount(c *fiber.Ctx) error {
	walletID := c.Params("id")
	walletUUID, err := uuid.Parse(walletID)
	if err != nil {
		h.logger.Error("GetWalletAmount Parse", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	amount, err := h.Uc.GetWalletAmount(c.Context(), walletUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		h.logger.Error("Error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info("Get wallet amount success", zap.Any("amount", amount))
	return c.Status(fiber.StatusOK).JSON(domain.GetWalletAmountResponse(amount))
}
