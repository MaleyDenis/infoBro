package reddit

import (
	"context"
	"testing"
	"time"

	"github.com/dzianismalei/infoBro/internal/config"
	"github.com/dzianismalei/infoBro/internal/models"
)

// MockStateRepository is a mock implementation of ChannelStateRepository
type MockStateRepository struct {
	states map[string]*models.ChannelState
}

// NewMockStateRepository creates a new mock state repository
func NewMockStateRepository() *MockStateRepository {
	return &MockStateRepository{
		states: make(map[string]*models.ChannelState),
	}
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

func TestChannelIDFromName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Subreddit name",
			input:    "golang",
			expected: "reddit:golang",
		},
		{
			name:     "Empty subreddit name",
			input:    "",
			expected: "reddit:",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ChannelIDFromName(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		config        config.RedditConfig
		expectError   bool
	}{
		{
			name: "Valid config",
			config: config.RedditConfig{
				Enabled: true,
				Subreddits: []config.SubredditConfig{
					{
						Name: "golang",
						URL:  "https://www.reddit.com/r/golang",
					},
				},
				Settings: config.RedditSettings{
					Timeout:      30 * time.Second,
					UserAgent:    "TestAgent",
					ClientID:     "test_client_id",
					ClientSecret: "test_client_secret",
					Username:     "test_username",
					Password:     "test_password",
					Limit:        25,
					Sort:         "hot",
				},
			},
			expectError: false,
		},
		{
			name: "Disabled connector",
			config: config.RedditConfig{
				Enabled: false,
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := NewMockStateRepository()
			_, err := New(test.config, mockRepo)
			
			if test.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			
			if !test.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}