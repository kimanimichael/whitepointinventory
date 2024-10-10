package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mike-kimani/whitepointinventory/internal/adapters/database"
	"github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
	"github.com/mike-kimani/whitepointinventory/internal/api/httpapi"
	"github.com/mike-kimani/whitepointinventory/internal/app"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Welcome to WhitePointInventory V2!")

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("couldn't find a port in this environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in this environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database: ", dbURL)

	db := sqlcdatabase.New(conn)
	// Init repositories
	userRepositorySQL := database.NewUserRepositorySQL(db)
	farmerRepositorySQl := database.NewFarmerRepositorySQL(db)
	purchasesRepositorySQL := database.NewPurchasesRepositorySQL(db)
	paymentsRepositorySQL := database.NewPaymentsRepositorySQL(db)

	//Init services
	userService := app.NewUserService(userRepositorySQL)
	farmerService := app.NewFarmerService(farmerRepositorySQl)
	purchasesService := app.NewPurchaseService(purchasesRepositorySQL)
	paymentsService := app.NewPaymentsService(paymentsRepositorySQL)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	//Init http handlers
	userHandler := httpapi.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)
	farmerHandler := httpapi.NewFarmerHandler(farmerService)
	farmerHandler.RegisterRoutes(router)
	purchasesHandler := httpapi.NewPurchasesHandler(purchasesService, userService)
	purchasesHandler.RegisterRoutes(router)
	paymentsHandler := httpapi.NewPaymentsHandler(paymentsService, userService)
	paymentsHandler.RegisterRoutes(router)

	actualRouter := chi.NewRouter()
	actualRouter.Mount("/whitepoint", router)

	srv := &http.Server{
		Handler: actualRouter,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on port %s", portString)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
