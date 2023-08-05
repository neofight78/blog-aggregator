package main

import (
	"com.github/neofight78/blog-aggregator/internal/database"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func errHandler(w http.ResponseWriter, _ *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func readiness(w http.ResponseWriter, _ *http.Request) {
	respondWithJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	})
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := struct {
		Name string `json:"name"`
	}{}
	err := decoder.Decode(&request)
	if err != nil {
		log.Printf("unable to decode request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid User Format")
		return
	}

	user := database.CreateUserParams{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      request.Name,
	}

	err = cfg.Queries.CreateUser(context.Background(), user)
	if err != nil {
		log.Printf("unable to create user: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func currentUser(w http.ResponseWriter, _ *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var params struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil || len(params.Name) == 0 || len(params.Url) == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid parameters")
		return
	}

	feed := database.CreateFeedParams{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID.String(),
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID.String(),
	}

	tx, err := cfg.DB.Begin()
	if err != nil {
		log.Printf("unable to create transaction: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	queriesTx := cfg.Queries.WithTx(tx)

	err = queriesTx.CreateFeed(r.Context(), feed)
	if err != nil {
		_ = tx.Rollback()

		log.Printf("unable to create feed: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = queriesTx.CreateFeedFollow(r.Context(), feedFollow)
	if err != nil {
		_ = tx.Rollback()
		log.Printf("unable to create feed follow: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		log.Printf("unable to commit transaction: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Feed       database.CreateFeedParams       `json:"feed"`
		FeedFollow database.CreateFeedFollowParams `json:"feed_follow"`
	}{
		Feed:       feed,
		FeedFollow: feedFollow,
	})
}

func (cfg *apiConfig) listFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.Queries.ListFeeds(r.Context())
	if err != nil {
		log.Printf("unable to list feeds: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}

func (cfg *apiConfig) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	var params struct {
		FeedID string `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid parameters")
		return
	}

	if _, err = uuid.Parse(params.FeedID); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid parameters")
		return
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.FeedID,
		UserID:    user.ID.String(),
	}

	err = cfg.Queries.CreateFeedFollow(r.Context(), feedFollow)
	if err != nil {
		log.Printf("unable to create feed follow: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusOK, feedFollow)
}

func (cfg *apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID := chi.URLParam(r, "feedFollowID")

	if _, err := uuid.Parse(feedFollowID); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid parameters")
		return
	}

	feedFollow, err := cfg.Queries.GetFeedFollow(r.Context(), feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Feed not found")
		return
	}

	if feedFollow.UserID != user.ID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	err = cfg.Queries.DeleteFeedFollow(r.Context(), feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func (cfg *apiConfig) listFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.Queries.ListFeedFollows(r.Context(), user.ID.String())
	if err != nil {
		log.Printf("unable to list feeds follows: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	
	respondWithJSON(w, http.StatusOK, feedFollows)
}
