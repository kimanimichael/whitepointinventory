package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main()  {
	fmt.Println("Welcome to whitepoint invetory")

	godotenv.Load(".env")

	portstring := os.Getenv("PORT")

	if portstring == "" {
		log.Fatal("couldn't find a port in this environment")
	}

	router := chi.NewRouter()
	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerHealth)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portstring,
	}
	log.Printf("Server starting on port %v", portstring)

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}