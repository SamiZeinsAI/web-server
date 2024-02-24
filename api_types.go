package main

import "github.com/SamiZeinsAI/web-server/internal/database"

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
	apiKey         string
}
