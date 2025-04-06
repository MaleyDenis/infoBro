package reddit

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dzianismalei/infoBro/internal/config"
	"github.com/dzianismalei/infoBro/internal/models"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Connector implements NewsConnector for Reddit
type Connector struct {
	client          *reddit.Client
	subreddits      []config.SubredditConfig
	limit           int
	sort            string
	stateRepository models.ChannelStateRepository
}

// New creates a new Reddit connector
func New(cfg config.RedditConfig, stateRepo models.ChannelStateRepository) (*Connector, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("reddit connector is disabled in config")
	}

	var client *reddit.Client
	var err error

	httpClient := &http.Client{
		Timeout: cfg.Settings.Timeout,
	}

	// Check if credentials are provided and not placeholder values
	if cfg.Settings.ClientID != "your_client_id" &&
		cfg.Settings.ClientSecret != "your_client_secret" &&
		cfg.Settings.Username != "your_username" &&
		cfg.Settings.Password != "your_password" {

		// Log all the variables
		log.Printf("ClientID: %s", cfg.Settings.ClientID)
		log.Printf("ClientSecret: %s", cfg.Settings.ClientSecret)
		log.Printf("Username: %s", cfg.Settings.Username)
		log.Printf("Password: %s", cfg.Settings.Password)

		// Use authenticated client
		credentials := reddit.Credentials{
			ID:       cfg.Settings.ClientID,
			Secret:   cfg.Settings.ClientSecret,
			Username: cfg.Settings.Username,
			Password: cfg.Settings.Password,
		}

		client, err = reddit.NewClient(credentials,
			reddit.WithHTTPClient(httpClient),
			reddit.WithUserAgent(cfg.Settings.UserAgent))

		if err != nil {
			return nil, fmt.Errorf("failed to create authenticated Reddit client: %w", err)
		}
	} else {
		// Use read-only client without authentication
		fmt.Println("Using read-only Reddit client. Rate limits will be lower.")
		userAgent := cfg.Settings.UserAgent
		if userAgent == "" || userAgent == "NewsAggregator/1.0 (by /u/your_username)" {
			userAgent = "Mozilla/5.0 (compatible; NewsAggregator/1.0)"
		}

		client = &reddit.Client{}
	}

	return &Connector{
		client:          client,
		subreddits:      cfg.Subreddits,
		limit:           cfg.Settings.Limit,
		sort:            cfg.Settings.Sort,
		stateRepository: stateRepo,
	}, nil
}

// GetNews retrieves news from Reddit
func (c *Connector) GetNews(ctx context.Context) ([]models.RawNews, error) {
	var allNews []models.RawNews

	for _, subreddit := range c.subreddits {
		// Set options for fetching posts
		listOptions := &reddit.ListPostOptions{
			ListOptions: reddit.ListOptions{
				Limit: c.limit,
			},
			Time: "day", // Default to day, can be customized based on c.sort if needed
		}

		// Fetch only top posts
		posts, _, err := c.client.Subreddit.TopPosts(ctx, subreddit.Name, listOptions)

		if err != nil {
			return nil, fmt.Errorf("failed to fetch top posts from r/%s: %w", subreddit.Name, err)
		}

		// Current time for FetchedAt field
		fetchedAt := time.Now()

		// Convert posts to RawNews
		for _, post := range posts {
			if post == nil || post.Created == nil {
				continue // Skip posts with missing data
			}

			// Create a new RawNews item
			newsItem := models.RawNews{
				SourceType:  "reddit",
				SourceID:    post.ID,
				SourceName:  subreddit.Name,
				SourceURL:   subreddit.URL, // Using URL from config
				Title:       post.Title,
				Content:     "", // Оставляем поле Content пустым
				URL:         fmt.Sprintf("https://www.reddit.com%s", post.Permalink),
				PublishedAt: post.Created.Time,
				FetchedAt:   fetchedAt,
				Metadata: map[string]interface{}{
					"author":           post.Author,
					"score":            post.Score,
					"numberOfComments": post.NumberOfComments,
					"isNSFW":           post.NSFW,
					"upvoteRatio":      post.UpvoteRatio,
					"subreddit":        subreddit.Name,
				},
			}

			// Добавляем все посты в список новостей
			allNews = append(allNews, newsItem)
		}
	}

	return allNews, nil
}
