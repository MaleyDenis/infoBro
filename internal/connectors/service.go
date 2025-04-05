package connectors

import (
	"context"
	"fmt"
	"sync"

	"github.com/dzianismalei/infoBro/internal/models"
)

// ConnectorService manages running connectors and storing their results
type ConnectorService struct {
	connectors map[string]models.NewsConnector
	storage    models.NewsStorage
	queue      models.NewsQueue
}

// NewConnectorService creates a new connector service
func NewConnectorService(connectors map[string]models.NewsConnector, storage models.NewsStorage, queue models.NewsQueue) *ConnectorService {
	return &ConnectorService{
		connectors: connectors,
		storage:    storage,
		queue:      queue,
	}
}

// RunConnector runs a specific connector and processes its results
func (s *ConnectorService) RunConnector(ctx context.Context, name string) (int, error) {
	connector, exists := s.connectors[name]
	if !exists {
		return 0, fmt.Errorf("connector %s not found", name)
	}

	// Get news from the connector
	news, err := connector.GetNews(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get news from %s: %w", name, err)
	}

	if len(news) == 0 {
		// No new items, but not an error
		return 0, nil
	}

	// Save the news to storage
	ids, err := s.storage.SaveRawNews(ctx, news)
	if err != nil {
		return 0, fmt.Errorf("failed to save news from %s: %w", name, err)
	}

	// Add the IDs to the queue for processing
	err = s.queue.AddToQueue(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to add news to queue from %s: %w", name, err)
	}

	return len(news), nil
}

// RunAllConnectors runs all available connectors in parallel
func (s *ConnectorService) RunAllConnectors(ctx context.Context) (map[string]ConnectorResult, error) {
	results := make(map[string]ConnectorResult)
	var wg sync.WaitGroup
	resultMutex := sync.Mutex{}

	for name := range s.connectors {
		wg.Add(1)
		go func(connectorName string) {
			defer wg.Done()

			count, err := s.RunConnector(ctx, connectorName)
			
			resultMutex.Lock()
			defer resultMutex.Unlock()
			
			if err != nil {
				results[connectorName] = ConnectorResult{
					Status:    "error",
					Message:   err.Error(),
					Processed: 0,
				}
			} else {
				results[connectorName] = ConnectorResult{
					Status:    "success",
					Processed: count,
				}
			}
		}(name)
	}

	wg.Wait()
	return results, nil
}

// ConnectorResult represents the result of running a connector
type ConnectorResult struct {
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
	Processed int    `json:"processed"`
}