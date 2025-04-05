package storage

import (
	"context"
	"errors"
	"time"

	"github.com/dzianismalei/infoBro/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB implements storage interfaces using MongoDB as a backend
type MongoDB struct {
	client             *mongo.Client
	database           string
	rawCollection      string
	processedCollection string
	channelStateCollection string
}

// NewMongoDB creates a new MongoDB storage instance
func NewMongoDB(uri, database, rawColl, processedColl, channelStateColl string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		client:                client,
		database:              database,
		rawCollection:         rawColl,
		processedCollection:   processedColl,
		channelStateCollection: channelStateColl,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// SaveRawNews saves raw news items to MongoDB and returns their ObjectIDs
func (m *MongoDB) SaveRawNews(ctx context.Context, news []models.RawNews) ([]primitive.ObjectID, error) {
	if len(news) == 0 {
		return []primitive.ObjectID{}, nil
	}

	collection := m.client.Database(m.database).Collection(m.rawCollection)
	
	var documents []interface{}
	for _, item := range news {
		doc := bson.M{
			"source_type":  item.SourceType,
			"source_id":    item.SourceID,
			"title":        item.Title,
			"content":      item.Content,
			"url":          item.URL,
			"published_at": item.PublishedAt,
			"fetched_at":   item.FetchedAt,
			"metadata":     item.Metadata,
		}
		documents = append(documents, doc)
	}

	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for _, id := range result.InsertedIDs {
		if oid, ok := id.(primitive.ObjectID); ok {
			ids = append(ids, oid)
		}
	}

	return ids, nil
}

// GetChannelState retrieves the state for a specific channel
func (m *MongoDB) GetChannelState(ctx context.Context, channelID string) (*models.ChannelState, error) {
	collection := m.client.Database(m.database).Collection(m.channelStateCollection)
	
	filter := bson.M{"channel_id": channelID}
	
	var result struct {
		ID                primitive.ObjectID `bson:"_id"`
		ChannelID         string             `bson:"channel_id"`
		LastMessageID     string             `bson:"last_message_id"`
		LastUpdateTime    time.Time          `bson:"last_update_time"`
		ProcessedMessages int                `bson:"processed_messages"`
	}
	
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Return an empty state for new channels
			return &models.ChannelState{
				ChannelID:         channelID,
				LastMessageID:     "",
				LastUpdateTime:    time.Time{},
				ProcessedMessages: 0,
			}, nil
		}
		return nil, err
	}
	
	return &models.ChannelState{
		ChannelID:         result.ChannelID,
		LastMessageID:     result.LastMessageID,
		LastUpdateTime:    result.LastUpdateTime,
		ProcessedMessages: result.ProcessedMessages,
	}, nil
}

// UpdateChannelState updates the state for a specific channel
func (m *MongoDB) UpdateChannelState(ctx context.Context, state *models.ChannelState) error {
	collection := m.client.Database(m.database).Collection(m.channelStateCollection)
	
	filter := bson.M{"channel_id": state.ChannelID}
	update := bson.M{
		"$set": bson.M{
			"last_message_id":     state.LastMessageID,
			"last_update_time":    state.LastUpdateTime,
			"processed_messages":  state.ProcessedMessages,
		},
	}
	
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// SaveProcessedNews saves a processed news item to MongoDB
func (m *MongoDB) SaveProcessedNews(ctx context.Context, news models.ProcessedNews) (primitive.ObjectID, error) {
	collection := m.client.Database(m.database).Collection(m.processedCollection)
	
	result, err := collection.InsertOne(ctx, news)
	if err != nil {
		return primitive.NilObjectID, err
	}
	
	return result.InsertedID.(primitive.ObjectID), nil
}