package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lyfoore/subscriptions_service/internal/db"
	"github.com/lyfoore/subscriptions_service/internal/service"
	"github.com/lyfoore/subscriptions_service/internal/transport/http"
	"os"
)

// @title Subscriptions CRUDL API
// @version 1.0
// @description This is a subscriptions CRUDL service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@crudl.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("POSTGRES_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	postgresDB, err := db.NewPostgresDB(connStr)

	subscriptionsService := service.NewService(postgresDB)
	handler := http.NewHandler(subscriptionsService)
	engine := http.NewRouter(handler)
	engine.SetupRouter()
	engine.Run(":8080")
}
