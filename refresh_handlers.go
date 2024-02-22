package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SamiZeinsAI/web-server/internal/auth"
)

func (cfg *apiConfig) PostRefreshHandler(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, "Bearer token not found")
		return
	}
	dbStructure, err := cfg.DB.LoadDB()
	if err != nil {
		RespondWithError(w, 500, "Error loading database")
		return
	}
	if _, ok := dbStructure.RevokedTokens[tokenString]; ok {
		RespondWithError(w, 401, "This token has been revoked")
		return
	}

	token, err := auth.ParseToken(tokenString, cfg.jwtSecret)
	if err != nil {
		RespondWithError(w, 401, "Error parsing token")
		return
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer != "chirpy-refresh" {
		RespondWithError(w, 401, "Issuer not refresh token shape")
		return
	}
	id, err := token.Claims.GetSubject()
	if err != nil {
		RespondWithError(w, 401, "Error getting subject header")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		RespondWithError(w, 401, "Id not valid integer in subject header")
		return
	}
	accessToken, err := auth.MakeToken(idInt, "chirpy-access", cfg.jwtSecret)
	if err != nil {
		RespondWithError(w, 500, "Error making access token")
		return
	}

	respBody := returnVals{
		Token: accessToken,
	}
	w.Header().Set("Content-Type", "applicaiton/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&respBody)
	if err != nil {
		RespondWithError(w, 500, "Error encoding response")
		return
	}
}
