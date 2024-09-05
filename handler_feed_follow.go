package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GabrielPereira187/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type FeedFollowRequest struct {
	FeedId string `json:"feed_id"`
}

func (cfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, token string) {
	var param FeedFollowRequest

	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		respondWithError(w, 403, "Error to Unmarshall")
		return 
	}

	user, err := cfg.DB.GetUser(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Error to get user")
		return
	}

	uuidString, err := uuid.Parse(param.FeedId)
	if err != nil {
		respondWithError(w, 403, "Error parsing UUID")
		return
	}

	feedFollowToCreate := database.FeedFollow{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID: uuidString,
		UserID: user.ID,
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams(feedFollowToCreate))
	if err != nil {
		respondWithError(w, 403, "Error to create feed follow")
		return
	}

	respondWithJSON(w, 201, feedFollow)

}

func (cfg *apiConfig) handlerGetFeedFollowByUser(w http.ResponseWriter, r *http.Request, token string) {
	user, err := cfg.DB.GetUser(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Error to get user")
		return
	}

	feedFollows, err := cfg.DB.GetAllFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 403, "Error to create feed follow")
		return
	}

	respondWithJSON(w, 201, feedFollows)
}

func (cfg *apiConfig) handlerDeleteFeedFollowByUserID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    feedFollowID := vars["feedFollowID"]

	uuidString, err := uuid.Parse(feedFollowID)
	if err != nil {
		respondWithError(w, 403, "Error parsing UUID")
		return
	}

	err = cfg.DB.DeleteFeedFollows(r.Context(), uuidString)
	if err != nil {
		respondWithError(w, 403, "Error to create feed follow")
		return
	}

	respondWithJSON(w, 201, "Deleted")
}