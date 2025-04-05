package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewsConnector - interface for all news sources
type NewsConnector interface {
	GetNews(ctx context.Context) ([]RawNews, error)
}

// RawNews - structure for storing news in a standard format
type RawNews struct {
	SourceType  string
	SourceID    string
	SourceName  string
	SourceURL   string
	Title       string
	Content     string
	URL         string
	PublishedAt time.Time
	FetchedAt   time.Time
	Metadata    map[string]interface{}
}

// ChannelState - source state structure
type ChannelState struct {
	ChannelID         string
	LastMessageID     string
	LastUpdateTime    time.Time
	ProcessedMessages int
}

// ChannelStateRepository - interface for storing source states
type ChannelStateRepository interface {
	GetChannelState(ctx context.Context, channelID string) (*ChannelState, error)
	UpdateChannelState(ctx context.Context, state *ChannelState) error
}

// NewsStorage - interface for news storage
type NewsStorage interface {
	SaveRawNews(ctx context.Context, news []RawNews) ([]primitive.ObjectID, error)
}

// NewsQueue - interface for news queue
type NewsQueue interface {
	AddToQueue(ctx context.Context, newsIDs []primitive.ObjectID) error
}

// ProcessedNews - structure for storing processed news
type ProcessedNews struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RawID       primitive.ObjectID `json:"raw_id" bson:"raw_id"`
	Title       string             `json:"title" bson:"title"`
	Content     string             `json:"content" bson:"content"`
	SourceType  string             `json:"source_type" bson:"source_type"`
	SourceID    string             `json:"source_id" bson:"source_id"`
	SourceName  string             `json:"source_name" bson:"source_name"`
	SourceURL   string             `json:"source_url" bson:"source_url"`
	URL         string             `json:"url" bson:"url"`
	PublishedAt time.Time          `json:"published_at" bson:"published_at"`
	ProcessedAt time.Time          `json:"processed_at" bson:"processed_at"`
}