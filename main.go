package main

import (
	"finances-backend/handlers"
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
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Open database
	storage := storage.NewPostgresStorage()
	err = storage.Open(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	// defer storage.Close()

	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

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
	categoryService := services.NewCategoryService(storage)
	operationService := services.NewOperationService(storage)
	searchService := services.NewSearchService(storage)

	handler := handlers.NewHandler(*authService, *userService, *walletService, *categoryService, *operationService, *searchService)

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
	autorizedGroup.Get("/wallet/:id", handler.GetWalletById)
	autorizedGroup.Patch("/wallet/:id", handler.UpdateWallet)

	// Categories
	autorizedGroup.Post("/category", handler.CreateCategory)
	autorizedGroup.Get("/category", handler.GetCategories)
	autorizedGroup.Get("/category/:id", handler.GetCategoryById)
	autorizedGroup.Patch("/category/:id", handler.UpdateCategory)
	autorizedGroup.Delete("/category/:id", handler.DeleteCategory)

	// Operations
	autorizedGroup.Post("/operation", handler.CreateOperation)
	autorizedGroup.Get("/operation", handler.GetOperations)
	autorizedGroup.Get("/operation/:id", handler.GetOperationById)
	autorizedGroup.Patch("/operation/:id", handler.UpdateOperation)
	autorizedGroup.Delete("/operation/:id", handler.DeleteOperation)

	// Search
	autorizedGroup.Get("/search/operation", handler.SearchOperations)

	// Start the server
	app.Listen(":3000")
}
