package main

import (
	"log"
	wallets_handler "tech_task/internal/handler/wallets"
	wallets_usecase "tech_task/internal/usecase/wallets"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	app := fiber.New()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot create zap logger: %v", err)
	}
	defer logger.Sync()

	//cfg, err := config.Load()
	//if err != nil {
	//	log.Fatalf("failed to load config: %v", err)
	//}
	//
	//database, err := db.Connect(cfg)
	//if err != nil {
	//	log.Fatalf("failed to connect to db: %v", err)
	//}

	//userRepo := wallets_repo.New(database)
	userUseCase := wallets_usecase.NewUseCase(logger /* userRepo */)
	userHandler := wallets_handler.NewHandler(logger, userUseCase)
	wallets_handler.RegisterRoutes(app, userHandler)

	//serverPort := os.Getenv("API_PORT")
	//if serverPort == "" {
	//	log.Fatal("environment variable SERVER_PORT isn't set")
	//}

	log.Fatal(app.Listen(":" + "8080"))
}
