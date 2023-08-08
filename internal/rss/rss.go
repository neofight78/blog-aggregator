package rss

import (
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Item struct {
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	Published   time.Time `json:"published"`
	FeedID      uuid.UUID `json:"feed_id"`
}

type rss struct {
	Channel struct {
		Item []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func FetchFeedItems(feedID uuid.UUID, url string) ([]Item, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch feed %s: %w", url, err)
	}

	var feed rss

	decoder := xml.NewDecoder(response.Body)
	err = decoder.Decode(&feed)
	if err != nil {
		return nil, fmt.Errorf("unable to decode feed: %w", err)
	}

	var items []Item

	for _, item := range feed.Channel.Item {
		published, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return nil, fmt.Errorf("unable to parse date '%s': %w", item.PubDate, err)
		}

		items = append(items, Item{
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			Published:   published,
			FeedID:      feedID,
		})
	}

	return items, nil
}
