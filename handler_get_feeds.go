package main

import "net/http"


func (cfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, 401, "Error")
		return
	}

	respondWithJSON(w, 200, feeds)
}