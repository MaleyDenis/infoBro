package queue

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RedisQueue implements NewsQueue using Redis lists
type RedisQueue struct {
	client *redis.Client
	queueKey string
	processingKey string
	failedKey string
}

// NewRedisQueue creates a new Redis queue
func NewRedisQueue(address, password string, db int) (*RedisQueue, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	// Check connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisQueue{
		client:        client,
		queueKey:      "news:queue",
		processingKey: "news:processing",
		failedKey:     "news:failed",
	}, nil
}

// Close closes the Redis connection
func (r *RedisQueue) Close() error {
	return r.client.Close()
}

// AddToQueue adds news IDs to the queue
func (r *RedisQueue) AddToQueue(ctx context.Context, newsIDs []primitive.ObjectID) error {
	if len(newsIDs) == 0 {
		return nil
	}

	// Convert ObjectIDs to string values
	values := make([]interface{}, len(newsIDs))
	for i, id := range newsIDs {
		values[i] = id.Hex()
	}

	// Add to the queue (RPUSH)
	return r.client.RPush(ctx, r.queueKey, values...).Err()
}

// GetFromQueue retrieves a news ID from the queue and moves it to the processing queue
func (r *RedisQueue) GetFromQueue(ctx context.Context) (string, error) {
	// Atomically move from queue to processing list (BLMOVE)
	result, err := r.client.BLMove(ctx, r.queueKey, r.processingKey, "LEFT", "RIGHT", 0).Result()
	if err != nil {
		return "", err
	}
	
	return result, nil
}

// AcknowledgeProcessed removes a processed item from the processing queue
func (r *RedisQueue) AcknowledgeProcessed(ctx context.Context, newsID string) error {
	// Remove specific item from processing list
	count, err := r.client.LRem(ctx, r.processingKey, 1, newsID).Result()
	if err != nil {
		return err
	}
	
	if count == 0 {
		return fmt.Errorf("item %s not found in processing queue", newsID)
	}
	
	return nil
}

// MarkAsFailed moves an item from processing to failed queue
func (r *RedisQueue) MarkAsFailed(ctx context.Context, newsID string) error {
	// First remove from processing
	count, err := r.client.LRem(ctx, r.processingKey, 1, newsID).Result()
	if err != nil {
		return err
	}
	
	if count == 0 {
		return fmt.Errorf("item %s not found in processing queue", newsID)
	}
	
	// Then add to failed
	return r.client.RPush(ctx, r.failedKey, newsID).Err()
}

// GetQueueStats returns count of items in each queue
func (r *RedisQueue) GetQueueStats(ctx context.Context) (map[string]int64, error) {
	pipe := r.client.Pipeline()
	
	queueCount := pipe.LLen(ctx, r.queueKey)
	processingCount := pipe.LLen(ctx, r.processingKey)
	failedCount := pipe.LLen(ctx, r.failedKey)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	
	return map[string]int64{
		"queue":      queueCount.Val(),
		"processing": processingCount.Val(),
		"failed":     failedCount.Val(),
	}, nil
}