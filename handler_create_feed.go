package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GabrielPereira187/blog-aggregator/internal/database"
	"github.com/google/uuid"
)


type parameters struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type FeedResponse struct {
	Feed database.Feed `json:"feed"`
	FeedFollow database.FeedFollow `json:"feed_follow"`
}



func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, token string) {
	var param parameters

	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		respondWithError(w, 401, "Error to unmarshall")
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Error to get user")
		return
	}

	feedToCreate := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url: param.Url,
		Name: param.Name,
		UserID: user.ID,
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), feedToCreate)
	if err != nil {
		respondWithError(w, 401, "Error to insert feed")
		return
	}

	feedFollowToCreate := database.FeedFollow{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID: feed.ID,
		UserID: user.ID,
	}


	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams(feedFollowToCreate)) 
	if err != nil {
		respondWithError(w, 401, "Error to insert feed_follow")
		return
	}

	respondWithJSON(w, 201, FeedResponse{
		Feed: feed,
		FeedFollow: feedFollow,
	})

}