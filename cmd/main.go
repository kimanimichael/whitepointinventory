package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
	"github.com/mike-kimani/whitepointinventory/internal/farmers"
	"github.com/mike-kimani/whitepointinventory/internal/farmers/api"
	"github.com/mike-kimani/whitepointinventory/internal/http"
	"github.com/mike-kimani/whitepointinventory/internal/payments"
	"github.com/mike-kimani/whitepointinventory/internal/payments/api"
	"github.com/mike-kimani/whitepointinventory/internal/purchases"
	"github.com/mike-kimani/whitepointinventory/internal/purchases/api"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"github.com/mike-kimani/whitepointinventory/internal/users/api"
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
	userRepositorySQL := users.NewUserRepositorySQL(db)
	farmerRepositorySQl := farmers.NewFarmerRepositorySQL(db)
	purchasesRepositorySQL := purchases.NewPurchaseRepositorySQL(db)
	paymentsRepositorySQL := payments.NewPaymentsRepositorySQL(db)

	//Init services
	userService := users.NewUserService(userRepositorySQL)
	farmerService := farmers.NewFarmerService(farmerRepositorySQl)
	purchasesService := purchases.NewPurchaseService(purchasesRepositorySQL)
	paymentsService := payments.NewPaymentsService(paymentsRepositorySQL)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	//Init http handlers
	userHandler := usersapi.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)
	farmerHandler := farmersapi.NewFarmerHandler(farmerService)
	farmerHandler.RegisterRoutes(router)
	purchasesHandler := purchasesapi.NewPurchasesHandler(purchasesService, userService)
	purchasesHandler.RegisterRoutes(router)
	paymentsHandler := paymentsapi.NewPaymentsHandler(paymentsService, userService)
	paymentsHandler.RegisterRoutes(router)

	httpapi.RegisterHealthHandlerRoutes(router)

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
