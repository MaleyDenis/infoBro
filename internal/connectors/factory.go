package connectors

import (
	"fmt"

	"github.com/dzianismalei/infoBro/internal/config"
	"github.com/dzianismalei/infoBro/internal/connectors/reddit"
	"github.com/dzianismalei/infoBro/internal/models"
)

// Factory creates connectors based on configuration
type Factory struct {
	config         *config.ConnectorsConfig
	stateRepository models.ChannelStateRepository
}

// NewFactory creates a new connector factory
func NewFactory(cfg *config.ConnectorsConfig, stateRepo models.ChannelStateRepository) *Factory {
	return &Factory{
		config:         cfg,
		stateRepository: stateRepo,
	}
}

// CreateRedditConnector creates a Reddit connector
func (f *Factory) CreateRedditConnector() (models.NewsConnector, error) {
	if !f.config.Reddit.Enabled {
		return nil, fmt.Errorf("reddit connector is disabled in config")
	}
	
	return reddit.New(f.config.Reddit, f.stateRepository)
}

// CreateAllConnectors creates all enabled connectors
func (f *Factory) CreateAllConnectors() (map[string]models.NewsConnector, error) {
	connectors := make(map[string]models.NewsConnector)
	
	// Create Reddit connector if enabled
	if f.config.Reddit.Enabled {
		connector, err := f.CreateRedditConnector()
		if err != nil {
			return connectors, fmt.Errorf("failed to create Reddit connector: %w", err)
		}
		connectors["reddit"] = connector
	}
	
	// Here you would add creation of other connector types (Telegram, RSS, etc.)
	
	return connectors, nil
}