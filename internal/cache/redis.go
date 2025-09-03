package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(key string, value any, expiration time.Duration) error
	Get(key string, dest any) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Flush() error
	// Additional methods for advanced caching
	SetWithTags(key string, value any, tags []string, expiration time.Duration) error
	InvalidateTag(tag string) error
	InvalidateTags(tags []string) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (r *RedisCache) Set(key string, value any, expiration time.Duration) error {
	ctx := context.Background()

	// Convert value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisCache) Get(key string, dest any) error {
	ctx := context.Background()

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (r *RedisCache) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Exists(key string) (bool, error) {
	ctx := context.Background()
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *RedisCache) Flush() error {
	ctx := context.Background()
	return r.client.FlushDB(ctx).Err()
}

// Advanced caching with tags
func (r *RedisCache) SetWithTags(key string, value any, tags []string, expiration time.Duration) error {
	ctx := context.Background()

	// Convert value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Set the main key
	pipe := r.client.TxPipeline()
	pipe.Set(ctx, key, data, expiration)

	// Add key to tag sets
	for _, tag := range tags {
		tagKey := fmt.Sprintf("tag:%s", tag)
		pipe.SAdd(ctx, tagKey, key)
		// Set expiration for tag key (only if it doesn't exist to avoid overriding)
		pipe.Expire(ctx, tagKey, expiration)
	}

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisCache) InvalidateTag(tag string) error {
	return r.InvalidateTags([]string{tag})
}

func (r *RedisCache) InvalidateTags(tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	ctx := context.Background()

	// Collect all keys to delete
	var allKeys []string

	// Get keys for each tag
	for _, tag := range tags {
		tagKey := fmt.Sprintf("tag:%s", tag)

		// Get all keys associated with this tag
		keys, err := r.client.SMembers(ctx, tagKey).Result()
		if err != nil {
			// Continue with other tags if one fails
			continue
		}

		allKeys = append(allKeys, keys...)
		// Also delete the tag set itself
		allKeys = append(allKeys, tagKey)
	}

	// Delete all collected keys
	if len(allKeys) > 0 {
		return r.client.Del(ctx, allKeys...).Err()
	}

	return nil
}

// Cache keys helpers
const (
	UserCachePrefix   = "user:"
	RoleCachePrefix   = "role:"
	UsersListTag      = "users:list"
	UsersTag          = "users"
	DefaultExpiration = 1 * time.Hour
	LongExpiration    = 24 * time.Hour
	ShortExpiration   = 10 * time.Minute
)
