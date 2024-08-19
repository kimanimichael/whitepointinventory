package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mike-kimani/whitepointinventory/internal/database"
	"github.com/mike-kimani/whitepointinventory/internal/farmers"
	"github.com/mike-kimani/whitepointinventory/internal/middleware"
	"github.com/mike-kimani/whitepointinventory/internal/payments"
	"github.com/mike-kimani/whitepointinventory/internal/purchases"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"github.com/mike-kimani/whitepointinventory/pkg/health"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Welcome to WhitePointInventory V2!")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
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

	db := database.New(conn)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	usersApiCfg := users.ApiConfig{
		DB: db,
	}

	purchasesApiCfg := purchases.ApiConfig{
		DB: db,
	}

	farmersApiCfg := farmers.ApiConfig{
		DB: db,
	}

	middlewareApiCfg := middleware.ApiConfig{
		DB: db,
	}

	paymentsApiCfg := payments.ApiConfig{
		DB: db,
	}

	v2Router := chi.NewRouter()
	v2Router.Get("/health", health.HandlerHealth)
	v2Router.Post("/users", usersApiCfg.HandlerCreateUser)
	v2Router.Get("/users", usersApiCfg.HandlerGetUserFromCookie)
	v2Router.Get("/user", usersApiCfg.HandlerGetUsers)
	v2Router.Post("/login", usersApiCfg.HandlerUserLogin)
	v2Router.Post("/logout", usersApiCfg.HandlerUserLogout)
	v2Router.Post("/farmers", farmersApiCfg.HandlerCreateFarmer)
	v2Router.Get("/farmers", farmersApiCfg.HandlerGetFarmerByName)
	v2Router.Get("/farmer", farmersApiCfg.HandlerGetFarmers)
	v2Router.Delete("/farmers/{farmer_id}", farmersApiCfg.HandlerDeleteFarmer)
	v2Router.Post("/purchases", middlewareApiCfg.MiddlewareAuth(purchasesApiCfg.HandlerCreatePurchases))
	v2Router.Get("/purchases", purchasesApiCfg.HandlerGetPurchases)
	v2Router.Get("/purchase", purchasesApiCfg.HandlerGetPurchaseByID)
	v2Router.Delete("/purchases/{purchase_id}", middlewareApiCfg.MiddlewareAuth(purchasesApiCfg.HandlerDeletePurchase))
	v2Router.Post("/payments", middlewareApiCfg.MiddlewareAuth(paymentsApiCfg.HandlerCreatePayment))
	v2Router.Get("/payment", paymentsApiCfg.HandlerGetPaymentByID)
	v2Router.Get("/payments", paymentsApiCfg.HandlerGetPayments)
	v2Router.Delete("/payments/{payment_id}", middlewareApiCfg.MiddlewareAuth(paymentsApiCfg.HandlerDeletePayment))

	router.Mount("/v2", v2Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on port %s", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
