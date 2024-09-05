package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetUserByApiKey(w http.ResponseWriter, r *http.Request, token string) {
	user, err := cfg.DB.GetUser(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Error")
		return
	}

	respondWithJSON(w, 200, user)
}