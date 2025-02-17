package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sebsvt/financial-service/handler"
	"github.com/sebsvt/financial-service/repository"
	"github.com/sebsvt/financial-service/service"
)

func main() {
	godotenv.Load()
	db := initDB()
	app := fiber.New()
	app.Use(cors.New())
	api := app.Group("/api")
	financial_account_repo := repository.NewFinancialAccountRepositoryDB(db)
	financial_account_srv := service.NewFinancialAccountService(financial_account_repo)
	financial_account_handler := handler.NewFinancialAccountHandler(financial_account_srv)
	payment_srv := service.NewPaymentService()
	payment_hanler := handler.NewPaymentHandler(payment_srv)

	api.Get("/financial_account", financial_account_handler.GetFinancialFromOrganisation)
	api.Post("/financial_account/setup-account/:organisation_id", financial_account_handler.SetupAccountHandler)
	api.Post("/payment/payment-intent", payment_hanler.CreatePaymentIntent)

	app.Listen(":8084")
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
