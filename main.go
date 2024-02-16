package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	apiCfg := apiConfig{}
	const port = "8080"
	filepathRoot := "."
	r := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)
	fsHandler = apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./assets/logo.png"))))
	r.Handle("/assets/logo.png", fsHandler)

	apiRouter.Post("/validate_chirp", validateChirpHandler)
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.resetMetrics)
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
