package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/mike-kimani/whitepointinventory/internal/database"
	
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main()  {
	fmt.Println("Welcome to whitepoint invetory")

	godotenv.Load(".env")

	portstring := os.Getenv("PORT")

	if portstring == "" {
		log.Fatal("couldn't find a port in this environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in this environment")
	}

	conn, err := sql.Open("postgres", dbURL)	
	if err != nil {
		fmt.Println("Error: ", err)
		log.Fatal("Cannot connect to database")
	}

	db := database.New(conn)

	apiCfg := apiConfig {
		DB: db,
	}

	router := chi.NewRouter()
	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerHealth)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Post("/farmers", apiCfg.handlerCreateFarmer)
	v1Router.Get("/farmers", apiCfg.handlerGetFarmerByName)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portstring,
	}
	log.Printf("Server starting on port %v", portstring)

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}