package main

import (
	"com.github/neofight78/blog-aggregator/internal/database"
	"com.github/neofight78/blog-aggregator/internal/rss"
	"context"
	"github.com/google/uuid"
	"log"
	"strings"
	"sync"
	"time"
)

func fetchFeeds(cfg apiConfig) {
	const n = 10

	for {
		feeds, err := cfg.Queries.GetNextFeedsToFetch(context.Background(), n)

		if err == nil {
			results := make(chan []rss.Item)
			var waitGroup sync.WaitGroup

			for _, feed := range feeds {
				waitGroup.Add(1)
				feed := feed
				go func() {
					log.Printf("Fetching: %s\n", feed.Name)
					items, err := rss.FetchFeedItems(feed.ID, feed.Url)
					if err != nil {
						log.Printf("unable to fetch posts: %v", err)
						results <- []rss.Item{}
					} else {
						results <- items
					}
					waitGroup.Done()
				}()
			}

			go func() {
				waitGroup.Wait()
				close(results)
			}()

			log.Printf("Processing results\n")

			for result := range results {
				log.Printf("Got %d\n", len(result))
				for _, item := range result {
					err = cfg.Queries.CreatePost(context.Background(), database.CreatePostParams{
						ID:          uuid.NewString(),
						CreatedAt:   time.Now().UTC(),
						UpdatedAt:   time.Now().UTC(),
						Title:       item.Title,
						Url:         item.Url,
						Description: item.Description,
						PublishedAt: item.Published,
						FeedID:      item.FeedID.String(),
					})
					if err != nil && !strings.Contains(err.Error(), "Duplicate entry") {
						log.Printf("error saving post: %v", err)
					}
				}
			}
		} else {
			log.Fatalf("unable to fetch feeds: %v", err)
		}

		log.Printf("Results processed\n")

		time.Sleep(time.Minute)
	}
}
