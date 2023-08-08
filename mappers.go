package main

import (
	"com.github/neofight78/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

func mapCreateFeedParam(feed database.CreateFeedParams) Feed {
	var lastFetchedAt *time.Time

	return Feed{
		ID:            uuid.MustParse(feed.ID),
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        uuid.MustParse(feed.UserID),
		LastFetchedAt: lastFetchedAt,
	}
}
