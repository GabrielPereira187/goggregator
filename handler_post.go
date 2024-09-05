package main

import (
	"net/http"
	"strconv"

	"github.com/GabrielPereira187/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, nil)
}

func (cfg *apiConfig) handlerGetPostByUser(w http.ResponseWriter, r *http.Request, token string) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		respondWithError(w, 500, "erro")
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), token)
	if err != nil {
		respondWithError(w, 500, "erro")
		return
	}

	posts , err := cfg.DB.GetAllRecentPostsByUser(r.Context(), database.GetAllRecentPostsByUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		respondWithError(w, 500, "erro")
		return
	}

	respondWithJSON(w, 200, posts)
}