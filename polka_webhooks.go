package main

import (
	"encoding/json"
	"net/http"

	"github.com/SamiZeinsAI/web-server/internal/auth"
	"github.com/SamiZeinsAI/web-server/internal/database"
)

func (cfg *apiConfig) PostPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string         `json:"event"`
		Data  map[string]int `json:"data"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, 400, "Error decoding request body")
		return
	}
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil || cfg.apiKey != apiKey {
		RespondWithError(w, 401, "correct api key not included")
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}
	userID, ok := params.Data["user_id"]
	if !ok {
		RespondWithError(w, 400, "data field improper shape")
		return
	}
	dbStructure, err := cfg.DB.LoadDB()
	if err != nil {
		RespondWithError(w, 500, "Error loading database")
		return
	}
	user, ok := dbStructure.Users[userID]
	if !ok {
		RespondWithError(w, 404, "This user ID does not exist")
		return
	}
	dbStructure.Users[userID] = database.User{
		Id:          user.Id,
		Email:       user.Email,
		Password:    user.Password,
		IsChirpyRed: true,
	}
	err = cfg.DB.WriteDB(dbStructure)
	if err != nil {
		RespondWithError(w, 500, "Error writing database to json file")
		return
	}
	w.WriteHeader(http.StatusOK)
}
