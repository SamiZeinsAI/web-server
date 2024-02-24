package main

import (
	"flag"
	"log"
	"net/http"
	"os"

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
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.WriteDB(database.DBStructure{
			Chirps:        map[int]database.Chirp{},
			Users:         map[int]database.User{},
			RevokedTokens: map[string]database.RevokedToken{},
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
		apiKey:         os.Getenv("POLKA_API_KEY"),
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

	apiRouter.Get("/chirps", apiCfg.GetChirpsHandler)
	apiRouter.Post("/chirps", apiCfg.PostChirpHandler)
	apiRouter.Get("/chirps/{id}", apiCfg.GetChirpHandler)
	apiRouter.Delete("/chirps/{chirpID}", apiCfg.DeleteChirpHandler)

	apiRouter.Post("/users", apiCfg.PostUserHandler)
	apiRouter.Put("/users", apiCfg.PutUserHandler)

	apiRouter.Post("/login", apiCfg.PostLoginHandler)

	apiRouter.Post("/refresh", apiCfg.PostRefreshHandler)

	apiRouter.Post("/revoke", apiCfg.PostRevokeHandler)

	apiRouter.Post("/polka/webhooks", apiCfg.PostPolkaWebhook)

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
