package data

import (
	"context"
	"testing"
	"time"

	"github.com/kordlab/marketplace/config"
	"github.com/stretchr/testify/assert"
)

func TestNewMongoDB(t *testing.T) {
	cfg := config.LoadConfig()

	mongo, err := NewMongoDB(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, mongo)

	ctx := context.Background()
	defer mongo.Close(ctx)

	// Test collection helpers
	assert.NotNil(t, mongo.Users())
	assert.NotNil(t, mongo.Products())
	assert.NotNil(t, mongo.Purchases())
	assert.NotNil(t, mongo.Client())

	// Test JWT secret
	assert.Equal(t, cfg.JWTSecret, mongo.GetJWTSecret())

	// Test connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err = NewMongoDB(&config.Config{
		MongoURL:     "mongodb://invalid:27017",
		DatabaseName: "test_db",
	})
	assert.Error(t, err)
}

func TestCreateIndexes(t *testing.T) {
	cfg := config.LoadConfig()

	mongo, err := NewMongoDB(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, mongo)

	ctx := context.Background()
	defer mongo.Close(ctx)

	// Test that indexes were created
	indexes, err := mongo.Users().Indexes().List(ctx)
	assert.NoError(t, err)

	var count int
	for indexes.Next(ctx) {
		count++
	}
	// Count should be 3 (2 created + 1 default _id)
	assert.Equal(t, 3, count)
}
