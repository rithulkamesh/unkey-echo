package data

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kordlab/marketplace/config"
)

type RedisDB struct {
	client *redis.Client
	config *config.Config
}

func NewRedisDB(cfg *config.Config) (*RedisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisDB{
		client: client,
		config: cfg,
	}, nil
}

// Session Management
func (r *RedisDB) StoreSession(userID string, token string, expiry time.Duration) error {
	return r.client.Set(context.Background(), "session:"+token, userID, expiry).Err()
}

// Token Blacklisting
func (r *RedisDB) BlacklistToken(token string, expiry time.Duration) error {
	return r.client.Set(context.Background(), "blacklist:"+token, true, expiry).Err()
}

func (r *RedisDB) IsTokenBlacklisted(token string) bool {
	exists, _ := r.client.Exists(context.Background(), "blacklist:"+token).Result()
	return exists > 0
}

// Rate Limiting
func (r *RedisDB) IncrementRequestCount(ip string) (int64, error) {
	key := "ratelimit:" + ip
	pipe := r.client.Pipeline()
	incr := pipe.Incr(context.Background(), key)
	pipe.Expire(context.Background(), key, time.Minute)
	_, err := pipe.Exec(context.Background())
	return incr.Val(), err
}

// Caching
func (r *RedisDB) CacheUser(user *User) error {
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(),
		"user:"+user.ID.Hex(),
		userData,
		time.Hour).Err()
}

func (r *RedisDB) GetCachedUser(userID string) (*User, error) {
	data, err := r.client.Get(context.Background(), "user:"+userID).Bytes()
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *RedisDB) Close() error {
	return r.client.Close()
}
