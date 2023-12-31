package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/skovranek/rss_aggregator/internal/database"
)

func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user User) {
	feedParams, err := getFeedParams(r.Body)
	if err != nil {
		log.Printf("Error: handlerFeedsCreate: getFeedParams(r.Body): %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to decode request body")
		return
	}

	name := feedParams.Name
	url := feedParams.URL
	id := uuid.New()
	now := time.Now().UTC()
	userID := user.ID
	ctx := context.Background()

	dbFeed, err := cfg.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    userID,
	})
	if err != nil {
		log.Printf("Error: handlerFeedsCreate: database.CreateFeed: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to add feed to database")
		return
	}

	feed := databaseFeedToFeed(dbFeed)

	followID := uuid.New()

	dbFollow, err := cfg.DB.CreateFollow(ctx, database.CreateFollowParams{
		ID:        followID,
		FeedID:    feed.ID,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		log.Printf("Error: handlerFeedsCreate: cfg.DB.CreateFollow: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to add follow to database")
		return
	}

	follow := databaseFollowToFollow(dbFollow)

	respBody := struct {
		Feed   Feed   `json:"feed"`
		Follow Follow `json:"feed_follow"`
	}{
		Feed:   feed,
		Follow: follow,
	}

	respondWithJSON(w, http.StatusCreated, respBody)
}
