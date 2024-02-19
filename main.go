package main

import (
	"log"
	"net/http"

	"github.com/SamiZeinsAI/web-server/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	const port = "8080"
	filepathRoot := "."

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	r := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)
	fsHandler = apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./assets/logo.png"))))
	r.Handle("/assets/logo.png", fsHandler)

	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.resetMetrics)
	apiRouter.Get("/chirps", apiCfg.DB.GetChirpsHandler)
	apiRouter.Post("/chirps", apiCfg.DB.PostChirpHandler)

	apiRouter.Get("/chirps/{id}", apiCfg.DB.GetChirpHandler)
	apiRouter.Post("/users", apiCfg.DB.PostUserHandler)
	apiRouter.Post("/login", apiCfg.DB.PostUserLoginHandler)

	adminRouter.Get("/metrics", apiCfg.getMetrics)
	r.Mount("/api", apiRouter)
	r.Mount("/admin", adminRouter)
	corsR := middlewareCors(r)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsR,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
