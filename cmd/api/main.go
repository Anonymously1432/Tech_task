package main

import (
	"log"
	"os"
	"tech_task/db"
	"tech_task/internal/config"
	wallets_handler "tech_task/internal/handler/wallets"
	"tech_task/internal/helper"
	wallets_repo "tech_task/internal/repository/wallets"
	wallets_usecase "tech_task/internal/usecase/wallets"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	app := fiber.New()

	//app.Use(timeout.NewWithContext(
	//	func(c *fiber.Ctx) error {
	//		return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
	//			"error": "request timed out",
	//		})
	//	}, 30*time.Second, fiber.ErrRequestTimeout))

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot create zap logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	walletLocker := helper.NewWalletLocker()

	userRepo := wallets_repo.New(database)
	userUseCase := wallets_usecase.NewUseCase(logger, userRepo, walletLocker)
	userHandler := wallets_handler.NewHandler(logger, userUseCase)
	wallets_handler.RegisterRoutes(app, userHandler)

	serverPort := os.Getenv("API_PORT")
	if serverPort == "" {
		log.Fatal("environment variable SERVER_PORT isn't set")
	}

	log.Fatal(app.Listen(":" + "8080"))
}
