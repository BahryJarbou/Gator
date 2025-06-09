package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/BahryJarbou/Gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %v <limit(default 2)>", cmd.Name)
	}
	var limit int32
	limit = 2
	if len(cmd.Args) == 1 {
		limitInt, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		limit = int32(limitInt)

	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})

	if err != nil {
		return err
	}

	for i, post := range posts {
		feed, err := s.db.GetFeedByID(context.Background(), post.FeedID)
		if err != nil {
			log.Printf("couldn't find the feed for post titled %v: %v", post.Title, err)
		}
		fmt.Printf("Post #%d from feed %v:\n", i, feed.Name)
		printPost(post)
	}

	return nil
}

func printPost(post database.Post) {
	fmt.Println("- Title: ", post.Title)
	if post.Description.Valid {
		fmt.Println("- description: ", post.Description.String)
	}
	fmt.Println("- URL: ", post.Url)
	if post.PublishedAt.Valid {
		fmt.Println("- published at: ", post.PublishedAt.Time.Format(time.RFC1123))
	}
	fmt.Println("=========================================================")
	fmt.Println()
}
