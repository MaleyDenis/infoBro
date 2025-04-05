package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dzianismalei/infoBro/internal/config"
	"github.com/dzianismalei/infoBro/internal/models"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Connector implements NewsConnector for Reddit
type Connector struct {
	client           *reddit.Client
	subreddits       []config.SubredditConfig
	limit            int
	sort             string
	stateRepository  models.ChannelStateRepository
	useUnauthenticated bool
}

// New creates a new Reddit connector
func New(cfg config.RedditConfig, stateRepo models.ChannelStateRepository) (*Connector, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("reddit connector is disabled in config")
	}

	var client *reddit.Client
	var err error
	useUnauthenticated := false
	
	httpClient := &http.Client{
		Timeout: cfg.Settings.Timeout,
	}

	// Check if credentials are provided and not placeholder values
	if cfg.Settings.ClientID != "your_client_id" && 
	   cfg.Settings.ClientSecret != "your_client_secret" &&
	   cfg.Settings.Username != "your_username" &&
	   cfg.Settings.Password != "your_password" {
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
		useUnauthenticated = true
	}

	return &Connector{
		client:           client,
		subreddits:       cfg.Subreddits,
		limit:            cfg.Settings.Limit,
		sort:             cfg.Settings.Sort,
		stateRepository:  stateRepo,
		useUnauthenticated: useUnauthenticated,
	}, nil
}

// GetNews retrieves news from Reddit
func (c *Connector) GetNews(ctx context.Context) ([]models.RawNews, error) {
	var allNews []models.RawNews

	for _, subreddit := range c.subreddits {
		news, err := c.getSubredditNews(ctx, subreddit)
		if err != nil {
			// Log the error but continue with other subreddits
			fmt.Printf("Error getting news from %s: %v\n", subreddit.Name, err)
			continue
		}
		allNews = append(allNews, news...)
	}

	return allNews, nil
}

// getSubredditNews retrieves news from a single subreddit
func (c *Connector) getSubredditNews(ctx context.Context, subreddit config.SubredditConfig) ([]models.RawNews, error) {
	// Get the last state for this subreddit
	channelID := fmt.Sprintf("reddit:%s", subreddit.Name)
	state, err := c.stateRepository.GetChannelState(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel state: %w", err)
	}

	// Determine which posts to fetch based on the last post ID
	var lastPostID string
	if state.LastMessageID != "" {
		lastPostID = state.LastMessageID
	}

	// Get posts from the subreddit
	var posts []*reddit.Post
	var listErr error

	if c.useUnauthenticated {
		// Use unauthenticated method
		posts, listErr = fetchRedditPostsUnauthenticated(ctx, subreddit.Name, c.sort, c.limit)
		if listErr != nil {
			return nil, fmt.Errorf("failed to fetch posts unauthenticated: %w", listErr)
		}
	} else {
		// Use the authenticated API client
		listOptions := &reddit.ListOptions{
			Limit: c.limit,
		}
		
		switch strings.ToLower(c.sort) {
		case "hot":
			posts, _, listErr = c.client.Subreddit.HotPosts(ctx, subreddit.Name, listOptions)
		case "new":
			posts, _, listErr = c.client.Subreddit.NewPosts(ctx, subreddit.Name, listOptions)
		case "top":
			posts, _, listErr = c.client.Subreddit.TopPosts(ctx, subreddit.Name, &reddit.ListPostOptions{
				ListOptions: *listOptions,
				Time:        "day",
			})
		case "rising":
			posts, _, listErr = c.client.Subreddit.RisingPosts(ctx, subreddit.Name, listOptions)
		default:
			posts, _, listErr = c.client.Subreddit.HotPosts(ctx, subreddit.Name, listOptions)
		}

		if listErr != nil {
			return nil, fmt.Errorf("failed to list posts: %w", listErr)
		}
	}

	// Filter posts based on the last post ID
	var newPosts []*reddit.Post
	if lastPostID == "" {
		// If no last post ID, take all posts
		newPosts = posts
	} else {
		// Otherwise, take posts until we hit the last post ID
		for _, post := range posts {
			if post.ID == lastPostID {
				break
			}
			newPosts = append(newPosts, post)
		}
	}

	// Transform posts to RawNews format
	var news []models.RawNews
	now := time.Now()

	for _, post := range newPosts {
		// Skip posts that have no content
		if post.Title == "" {
			continue
		}

		// Create metadata
		metadata := map[string]interface{}{
			"subreddit":    post.SubredditName,
			"score":        post.Score,
			"upvote_ratio": post.UpvoteRatio,
			"num_comments": post.NumberOfComments,
			"permalink":    post.Permalink,
			"is_self":      post.IsSelfPost,
		}

		// Extract content from either selftext or url
		content := post.Body
		if content == "" && !post.IsSelfPost {
			content = fmt.Sprintf("External link: %s", post.URL)
		}

		// Create post URL
		postURL := fmt.Sprintf("https://www.reddit.com%s", post.Permalink)

		// Convert created time - use current time as fallback
		createdTime := now
		if post.Created != nil {
			createdTime = post.Created.Time
		}

		// Create RawNews item
		newsItem := models.RawNews{
			SourceType:  "reddit",
			SourceID:    subreddit.Name,
			SourceName:  fmt.Sprintf("r/%s", subreddit.Name),
			SourceURL:   subreddit.URL,
			Title:       post.Title,
			Content:     content,
			URL:         postURL,
			PublishedAt: createdTime,
			FetchedAt:   now,
			Metadata:    metadata,
		}

		news = append(news, newsItem)
	}

	// Update the channel state if we have new posts
	if len(newPosts) > 0 {
		state.LastMessageID = newPosts[0].ID
		state.LastUpdateTime = now
		state.ProcessedMessages += len(newPosts)

		err = c.stateRepository.UpdateChannelState(ctx, state)
		if err != nil {
			// Log error but continue
			fmt.Printf("Error updating channel state for %s: %v\n", channelID, err)
		}
	}

	return news, nil
}

// ChannelIDFromName generates a consistent channel ID for a subreddit
func ChannelIDFromName(name string) string {
	return fmt.Sprintf("reddit:%s", name)
}

// fetchRedditPostsUnauthenticated fetches posts directly from Reddit's JSON API without authentication
func fetchRedditPostsUnauthenticated(ctx context.Context, subreddit, sortType string, limit int) ([]*reddit.Post, error) {
	// Construct the URL to Reddit's JSON API
	url := fmt.Sprintf("https://www.reddit.com/r/%s/%s.json?limit=%d", 
		subreddit, strings.ToLower(sortType), limit)
	
	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set a user agent to avoid 429 errors
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; InfoBroNewsAggregator/1.0)")
	
	// Send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Parse the JSON response
	var redditResponse struct {
		Data struct {
			Children []struct {
				Data struct {
					ID            string  `json:"id"`
					Subreddit     string  `json:"subreddit"`
					Title         string  `json:"title"`
					Selftext      string  `json:"selftext"`
					URL           string  `json:"url"`
					Permalink     string  `json:"permalink"`
					Created       float64 `json:"created_utc"`
					Score         int     `json:"score"`
					UpvoteRatio   float64 `json:"upvote_ratio"`
					NumComments   int     `json:"num_comments"`
					IsSelf        bool    `json:"is_self"`
					Author        string  `json:"author"`
					SubredditName string  `json:"subreddit_name_prefixed"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &redditResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	// Convert to reddit.Post objects
	var posts []*reddit.Post
	for _, child := range redditResponse.Data.Children {
		data := child.Data
		createdTime := time.Unix(int64(data.Created), 0)
		timestamp := reddit.Timestamp{Time: createdTime}
		
		post := &reddit.Post{
			ID:              data.ID,
			SubredditName:   data.Subreddit,
			Title:           data.Title,
			Body:            data.Selftext,
			URL:             data.URL,
			Permalink:       data.Permalink,
			Created:         &timestamp,
			Score:           data.Score,
			UpvoteRatio:     float32(data.UpvoteRatio),
			NumberOfComments: data.NumComments,
			IsSelfPost:      data.IsSelf,
		}
		
		posts = append(posts, post)
	}
	
	return posts, nil
}