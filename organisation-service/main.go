package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sebsvt/organisation-service/handler"
	"github.com/sebsvt/organisation-service/middlewares"
	"github.com/sebsvt/organisation-service/repository"
	"github.com/sebsvt/organisation-service/services"
)

func main() {
	godotenv.Load()
	secret := os.Getenv("SECRET_KEY")
	db := initDB()
	org_repo := repository.NewOrganisationRepositoryDB(db)
	org_member_repo := repository.NewOrganisationMemberRepository(db)
	org_srv := services.NewOrganisationService(org_repo, org_member_repo)
	org_handler := handler.NewOrganisationHandler(org_srv)

	app := fiber.New()
	app.Use(cors.New())
	api := app.Group("/api")

	api.Post("/organisation", middlewares.AuthRequired(secret), org_handler.CreateNewOrganisation)
	api.Get("/organisation/my-organisations", middlewares.AuthRequired(secret), org_handler.GetAllOrganisationsFromUserID)
	api.Get("/organisation/:domain", org_handler.GetOrganisationFromDomain)
	api.Get("/organisation/:organisation_id/member", middlewares.AuthRequired(secret), org_handler.GetAllOrganisationMember)

	app.Listen(":8081")
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
