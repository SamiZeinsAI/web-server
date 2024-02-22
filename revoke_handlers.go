package main

import (
	"net/http"
	"time"

	"github.com/SamiZeinsAI/web-server/internal/auth"
	"github.com/SamiZeinsAI/web-server/internal/database"
)

func (cfg *apiConfig) PostRevokeHandler(w http.ResponseWriter, r *http.Request) {
	dbStructure, err := cfg.DB.LoadDB()
	if err != nil {
		RespondWithError(w, 500, "Error loading database")
		return
	}
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 400, "Error getting bearer token")
		return
	}
	token, err := auth.ParseToken(tokenString, cfg.jwtSecret)
	if err != nil {
		RespondWithError(w, 400, "Error parsing token")
		return
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer != "chirpy-refresh" {
		RespondWithError(w, 400, "Issuer in claims doesn't match the shape of a refresh token")
		return
	}

	dbStructure.RevokedTokens[tokenString] = database.RevokedToken{
		TimeRevoked: time.Now().UTC(),
		TokenString: tokenString,
	}
	cfg.DB.WriteDB(dbStructure)
	w.WriteHeader(http.StatusOK)

}
