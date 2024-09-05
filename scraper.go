package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/GabrielPereira187/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines...", timeBetweenRequest, concurrency)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <- ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))

		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue
		}
		log.Printf("Found %v feeds to fetch!", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.UpdateLastFetched(context.Background(), feed.ID)

	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		log.Println("Found post", item.Title)

		_ ,err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID: feed.ID,
			Title: item.Title,
			Url: item.Link,
			Description: stringToNullString(item.Description),
			PublishedAt: timeToNullTime(stringToTimePointer(item.PubDate)),
		})

		if err != nil {
			log.Println(err)
		}

		log.Println("inserted successfully")

	}

}

func stringToNullString(s string) sql.NullString {
    if s != "" {
        return sql.NullString{String: s, Valid: true}
    }
    return sql.NullString{Valid: false}
}

func timeToNullTime(t *time.Time) sql.NullTime {
    if t != nil {
        return sql.NullTime{Time: *t, Valid: true}
    }
    return sql.NullTime{Valid: false}
}

func stringToTimePointer(timeString string) (*time.Time) {
    // Define o formato da string de tempo
    layout := "2006-01-02 15:04:05"
    
    // Faz o parsing da string para o tipo time.Time
    t, err := time.Parse(layout, timeString)
    if err != nil {
        return nil
    }
    
    // Retorna o ponteiro para o valor de time.Time
    return &t
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}

	return &rssFeed, nil
}