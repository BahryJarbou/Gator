package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/BahryJarbou/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %v <time_between_reqs>", cmd.Name)
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("time_between_reqs should be in the form of #h, #m, #s, or a combination of that: %v", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("couldn't retieve the next feed to fetch:", err)
		return
	}
	fmt.Printf("fetching feeds for %v...\n", feed.Name)

	defer markAsFetched(s, feed.ID)

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Println("couldn't fetch feed: ", err)
		return
	}

	fmt.Printf("Title for feed %v are: \n", feed.Name)
	fmt.Println()
	for _, item := range feedData.Channel.Item {
		dateParsed, err := time.Parse(time.RFC1123Z, item.PubDate)
		publishDate := sql.NullTime{
			Time:  dateParsed,
			Valid: true,
		}
		if err != nil {
			publishDate.Valid = false
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: publishDate,
			FeedID:      feed.ID,
		})
		var err1 *pq.Error
		if errors.As(err, &err1) {
			if err1.Code == "23505" {
				log.Println("skipping duplicate post")
				continue
			} else {
				log.Printf("error writing the post to the database: %v", err)
				continue
			}
		} else if err != nil {
			log.Printf("error saving post: %v", err)
			continue
		}

		fmt.Printf("- %v\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
	fmt.Println("=============================================================")
	fmt.Println()
}

func markAsFetched(s *state, feed_id uuid.UUID) {
	err := s.db.MarkFeedFetched(context.Background(), feed_id)
	if err != nil {
		log.Println("couldn't mark feed as fetched: ", err)
		return
	}
}
