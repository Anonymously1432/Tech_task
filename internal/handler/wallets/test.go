package wallets

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *Handler) Test(c *fiber.Ctx) error {
	err := h.Uc.Test()
	if err != nil {
		h.logger.Error("Error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info("Success")
	return c.SendStatus(fiber.StatusOK)
}
