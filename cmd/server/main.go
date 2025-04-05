package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dzianismalei/infoBro/internal/api"
	"github.com/dzianismalei/infoBro/internal/config"
	"github.com/dzianismalei/infoBro/internal/connectors"
	"github.com/dzianismalei/infoBro/internal/connectors/reddit"
	"github.com/dzianismalei/infoBro/internal/models"
	"github.com/dzianismalei/infoBro/internal/queue"
	"github.com/dzianismalei/infoBro/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockStateRepository is a mock implementation of ChannelStateRepository for testing
type MockStateRepository struct {
	states map[string]*models.ChannelState
}

// GetChannelState retrieves a channel state
func (m *MockStateRepository) GetChannelState(ctx context.Context, channelID string) (*models.ChannelState, error) {
	state, ok := m.states[channelID]
	if !ok {
		return &models.ChannelState{
			ChannelID:         channelID,
			LastMessageID:     "",
			LastUpdateTime:    time.Time{},
			ProcessedMessages: 0,
		}, nil
	}
	return state, nil
}

// UpdateChannelState updates a channel state
func (m *MockStateRepository) UpdateChannelState(ctx context.Context, state *models.ChannelState) error {
	m.states[state.ChannelID] = state
	return nil
}

// MockNewsStorage is a mock implementation of the NewsStorage interface
type MockNewsStorage struct {
	connector    *connectors.ConnectorService
	redditConfig config.RedditConfig
}

// GetNewsList gets a list of news items with pagination
func (m *MockNewsStorage) GetNewsList(filters map[string]interface{}, page, pageSize int) (*api.NewsListResult, error) {
	// For a more realistic implementation, let's get actual news from Reddit
	ctx := context.Background()
	
	// Use the Reddit config from our struct
	redditCfg := m.redditConfig
	if !redditCfg.Enabled {
		// Fall back to default mock data if Reddit is disabled
		return m.getMockNewsList(page, pageSize)
	}
	
	// Create a mock repository
	mockRepo := &MockStateRepository{states: make(map[string]*models.ChannelState)}
	
	// Create the Reddit connector
	redditConnector, err := reddit.New(redditCfg, mockRepo)
	if err != nil {
		return m.getMockNewsList(page, pageSize)
	}
	
	// Get actual news
	news, err := redditConnector.GetNews(ctx)
	if err != nil || len(news) == 0 {
		return m.getMockNewsList(page, pageSize)
	}
	
	// Convert to API format
	var items []api.NewsItem
	for _, item := range news {
		// Create preview (truncated content)
		preview := item.Content
		if len(preview) > 150 {
			preview = preview[:150] + "..."
		}
		
		items = append(items, api.NewsItem{
			ID:            primitive.NewObjectID().Hex(),
			Title:         item.Title,
			ContentPreview: preview,
			Content:       item.Content,
			SourceType:    item.SourceType,
			SourceID:      item.SourceID,
			SourceName:    item.SourceName,
			SourceURL:     item.SourceURL,
			URL:           item.URL,
			PublishedAt:   item.PublishedAt,
			ProcessedAt:   item.FetchedAt,
		})
	}
	
	// Apply pagination
	totalItems := len(items)
	totalPages := (totalItems + pageSize - 1) / pageSize // Ceiling division
	
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= totalItems {
		items = []api.NewsItem{}
	} else if end > totalItems {
		items = items[start:totalItems]
	} else {
		items = items[start:end]
	}
	
	return &api.NewsListResult{
		Items: items,
		Pagination: api.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			TotalItems: totalItems,
		},
	}, nil
}

// getMockNewsList returns mock news data
func (m *MockNewsStorage) getMockNewsList(page, pageSize int) (*api.NewsListResult, error) {
	// Create a simple mock response
	items := []api.NewsItem{
		{
			ID:            "mock-news-1",
			Title:         "Go 1.21 Released",
			ContentPreview: "The Go development team has announced the release of Go 1.21...",
			SourceType:    "reddit",
			SourceID:      "golang",
			SourceName:    "r/golang",
			SourceURL:     "https://www.reddit.com/r/golang",
			URL:           "https://www.reddit.com/r/golang/comments/123456",
			PublishedAt:   time.Now().Add(-24 * time.Hour),
			ProcessedAt:   time.Now().Add(-23 * time.Hour),
		},
		{
			ID:            "mock-news-2",
			Title:         "Introduction to Concurrency in Go",
			ContentPreview: "Learn how to leverage Go's powerful concurrency features...",
			SourceType:    "reddit",
			SourceID:      "golang",
			SourceName:    "r/golang",
			SourceURL:     "https://www.reddit.com/r/golang",
			URL:           "https://www.reddit.com/r/golang/comments/123457",
			PublishedAt:   time.Now().Add(-48 * time.Hour),
			ProcessedAt:   time.Now().Add(-47 * time.Hour),
		},
	}

	// Calculate pagination
	return &api.NewsListResult{
		Items: items,
		Pagination: api.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: 1,
			TotalItems: len(items),
		},
	}, nil
}

// GetNewsById gets a specific news item by ID
func (m *MockNewsStorage) GetNewsById(id string) (*api.NewsItem, error) {
	// Try to convert the ID to an ObjectID
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		// If not a valid ObjectID, return mock data
		return m.getMockNewsItem(id), nil
	}
	
	// Get all news from our GetNewsList method
	result, err := m.GetNewsList(map[string]interface{}{}, 1, 100)
	if err != nil || result == nil || len(result.Items) == 0 {
		return m.getMockNewsItem(id), nil
	}
	
	// Try to find the item with the given ID
	for _, item := range result.Items {
		if item.ID == id {
			// Make a copy to return the full content
			fullItem := item
			fullItem.Content = item.ContentPreview
			if len(fullItem.Content) > 0 && fullItem.Content[len(fullItem.Content)-3:] == "..." {
				fullItem.Content = fullItem.Content[:len(fullItem.Content)-3]
			}
			return &fullItem, nil
		}
	}
	
	// If item not found, return mock data
	return m.getMockNewsItem(id), nil
}

// getMockNewsItem returns a mock news item with the given ID
func (m *MockNewsStorage) getMockNewsItem(id string) *api.NewsItem {
	return &api.NewsItem{
		ID:         id,
		Title:      "Go 1.21 Released",
		Content:    "The Go development team has announced the release of Go 1.21, featuring improved generics support, enhanced error handling, and better performance for concurrent operations.",
		SourceType: "reddit",
		SourceID:   "golang",
		SourceName: "r/golang",
		SourceURL:  "https://www.reddit.com/r/golang",
		URL:        "https://www.reddit.com/r/golang/comments/123456",
		PublishedAt: time.Now().Add(-24 * time.Hour),
		ProcessedAt: time.Now().Add(-23 * time.Hour),
	}
}

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config/connectors.yaml", "Path to connectors config file")
	mongoURI := flag.String("mongo-uri", "mongodb://localhost:27017", "MongoDB connection URI")
	mongoDatabase := flag.String("mongo-db", "infoBro", "MongoDB database name")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis server address")
	redisPassword := flag.String("redis-password", "", "Redis password")
	redisDB := flag.Int("redis-db", 0, "Redis database number")
	httpAddr := flag.String("http-addr", ":8080", "HTTP server address")
	runConnector := flag.String("run-connector", "", "Run a specific connector (reddit, telegram, rss)")
	flag.Parse()

	// Load configuration
	connectorsConfig, err := config.LoadConnectorsConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize MongoDB storage
	mongoStorage, err := storage.NewMongoDB(
		*mongoURI,
		*mongoDatabase,
		"raw_news",
		"processed_news",
		"channel_states",
	)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize Redis queue
	redisQueue, err := queue.NewRedisQueue(*redisAddr, *redisPassword, *redisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Create connector factory
	connectorFactory := connectors.NewFactory(connectorsConfig, mongoStorage)

	// Create all enabled connectors
	connectorMap, err := connectorFactory.CreateAllConnectors()
	if err != nil {
		log.Printf("Warning: some connectors could not be created: %v", err)
	}

	// Create connector service
	connectorService := connectors.NewConnectorService(connectorMap, mongoStorage, redisQueue)

	// Check if we should run a specific connector
	if *runConnector != "" {
		log.Printf("Running connector: %s", *runConnector)
		
		if *runConnector == "reddit" {
			// Special handling for testing Reddit connector
			log.Printf("Testing Reddit connector directly")
			
			// Create mock repository
			mockRepo := &MockStateRepository{states: make(map[string]*models.ChannelState)}
			
			// Load Reddit config
			redditCfg := connectorsConfig.Reddit
			
			// Create the connector
			redditConnector, err := reddit.New(redditCfg, mockRepo)
			if err != nil {
				log.Fatalf("Failed to create Reddit connector: %v", err)
			}
			
			// Get news
			news, err := redditConnector.GetNews(context.Background())
			if err != nil {
				log.Fatalf("Failed to get news from Reddit: %v", err)
			}
			
			// Print results
			log.Printf("Successfully retrieved %d news items from Reddit", len(news))
			for i, item := range news {
				if i >= 5 {
					log.Printf("... and %d more items", len(news)-5)
					break
				}
				log.Printf("- %s", item.Title)
			}
		} else {
			// Normal connector service run
			count, err := connectorService.RunConnector(context.Background(), *runConnector)
			if err != nil {
				log.Fatalf("Failed to run connector %s: %v", *runConnector, err)
			}
			log.Printf("Successfully processed %d items from %s", count, *runConnector)
		}
		return
	}

	// Create a mock news storage for testing
	mockNewsStorage := &MockNewsStorage{
		connector: connectorService,
		redditConfig: connectorsConfig.Reddit,
	}
	
	// Create API
	apiHandler := api.NewAPI(connectorService, mockNewsStorage)

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Register API routes
	apiHandler.RegisterRoutes(r)

	// Create HTTP server
	server := &http.Server{
		Addr:    *httpAddr,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("HTTP server listening on %s", *httpAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}