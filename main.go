package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/mike-kimani/whitepointinventory/internal/database"
	"github.com/yarlson/chistaticmiddleware/static"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func newTestRouter(route string, handler http.HandlerFunc) *chi.Mux {
	r := chi.NewRouter()
	r.Get(route, handler)

	// staticFileDirectory := http.Dir("./assets/")
	
	// staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))

	

	return r
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

	staticDir := "./assets"

	staticConfig := static.Config {
		Fs: os.DirFS(staticDir),
		Root: ".",
		FilePrefix: "/assets/",
		CacheDuration: 24 * time.Hour,
		Debug: true,
	}



	router := chi.NewRouter()
	// v1Router := chi.NewRouter()
	router.Use(static.Handler(staticConfig))

	v1Router := newTestRouter("/healthz", handlerHealth)

	// v1Router.Get("/healthz", handlerHealth)
	v1Router.Get("/test_static_file", handlerStaticFiles)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Post("/farmers", apiCfg.handlerCreateFarmer)
	v1Router.Get("/farmers", apiCfg.handlerGetFarmerByName)
	v1Router.Delete("/farmers/{farmer_id}", apiCfg.handlerDeleteFarmer)
	v1Router.Post("/purchases", apiCfg.middlewareAuth(apiCfg.handerCreatePurchases))
	v1Router.Get("/purchases", apiCfg.handlerGetPurchases)
	v1Router.Get("/purchase", apiCfg.handlerGetPurchaseByID)
	v1Router.Delete("/purchases/{purchase_id}", apiCfg.middlewareAuth(apiCfg.handlerDeletePurchase))
	v1Router.Post("/payments", apiCfg.middlewareAuth(apiCfg.handlerCreatePayment))
	v1Router.Get("/payment",apiCfg.handlerGetPaymentByID)
	v1Router.Get("/payments",apiCfg.handlerGetPayments)
	v1Router.Delete("/payments/{payment_id}", apiCfg.middlewareAuth(apiCfg.handlerDeletePayment))
	

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