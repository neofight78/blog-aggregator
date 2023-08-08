package rss

import (
	"github.com/google/uuid"
	"testing"
)

func TestFetchFeedItems(t *testing.T) {
	items, err := FetchFeedItems(uuid.New(), "https://wagslane.dev/index.xml")
	if err != nil {
		t.Errorf("failed to fetch items: %v", err)
	}

	if len(items) == 0 {
		t.Errorf("no items found")
	}

	t.Logf("Hello")

	for _, item := range items {
		t.Logf("%s", item.Title)
	}
}
