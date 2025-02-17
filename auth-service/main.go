package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sebsvt/cmu-contest-2024/auth-service/handlers"
	"github.com/sebsvt/cmu-contest-2024/auth-service/repository"
	"github.com/sebsvt/cmu-contest-2024/auth-service/services"
)

func main() {
	godotenv.Load()
	db := initDB()
	user_repo := repository.NewUserRepositoryDB(db)
	user_srv := services.NewUserService(user_repo)
	auth_srv := services.NewAuth(user_repo, []byte(os.Getenv("SECRET_KEY")), time.Hour*2, time.Hour*(24*7))
	user_handler := handlers.NewAuthHandler(user_srv, auth_srv)

	app := fiber.New()
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://localhost:3000, http://912e-182-52-165-224.ngrok-free.app",
	// 	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	// 	AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	// }))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "",
	}))
	api := app.Group("/api")

	api.Post("/auth/authorize", user_handler.Authorize)
	api.Post("/auth/sign-up", user_handler.SignUp)
	api.Post("/auth/sign-in", user_handler.SignIn)
	api.Post("/auth/refresh", user_handler.RefreshToken)

	app.Listen(":8080")
}

func initDB() *sqlx.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_DB"),
		os.Getenv("DATABASE_SSLMODE"),
	)
	fmt.Println(dsn)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return db
}
