package main

import (
	"finances-backend/services"
	"finances-backend/storage"
	"log"
	"os"
	"strconv"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environmental variables
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Open database
 	storage := storage.NewPostgresStorage()
	storage.Open(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	defer storage.Close()

	// Get JWT parameters
	expiration, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		panic(err)
	}
	secret := os.Getenv("JWT_SECRET")

	// Initialize services and a handler
	authService := services.NewAuthService(storage, secret, expiration)
	userService := services.NewUserService(storage)
	walletService := services.NewWalletService(storage)

	handler := NewHandler(*authService, *userService, *walletService)

	// Initialize Fiber
	app := fiber.New()
	app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))

	// Routes:
    // Public access
    publicGroup := app.Group("")
	publicGroup.Post("/register", handler.Register)
    publicGroup.Post("/login", handler.Login)

    // Authorized access
    autorizedGroup := app.Group("")
    autorizedGroup.Use(jwtware.New(jwtware.Config{
        SigningKey: jwtware.SigningKey{
            Key: []byte(os.Getenv("JWT_SECRET")),
        },
    }))
	autorizedGroup.Use(handler.AuthorizeMiddleware)
    autorizedGroup.Get("/profile", handler.Profile)
	// Wallets
	autorizedGroup.Post("/wallet", handler.CreateWallet)
	autorizedGroup.Get("/wallet", handler.GetWallets)
	autorizedGroup.Get("/wallet/:id", handler.GetWalletByID)

	// Start the server
	app.Listen(":3000")
}
