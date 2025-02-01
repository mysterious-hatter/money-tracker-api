package main

import (
	"finances-backend/services"
	"finances-backend/storage"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/contrib/jwt"
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

	// Initialize services and a handler
	userService := services.NewUserService(storage)
	handler := NewHandler(*userService)

	// Initialize Fiber
	app := fiber.New()
	app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))
    // Available for unauthorized users
    publicGroup := app.Group("")
	publicGroup.Post("/Register", handler.Register)
    publicGroup.Post("/login", handler.Login)

    // For private access
    autorizedGroup := app.Group("")
    autorizedGroup.Use(jwtware.New(jwtware.Config{
        SigningKey: jwtware.SigningKey{
            Key: []byte(os.Getenv("JWT_SECRET")),
        },
    }))
    autorizedGroup.Get("/profile", handler.Profile)

	// Routes
	app.Get("/user", handler.GetAllUsers)

	app.Listen(":3000")
}
