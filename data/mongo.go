package data

import (
	"context"
	"time"

	"github.com/kordlab/marketplace/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
	config   *config.Config
	redis    *RedisDB // Add Redis client
}

// Collections
const (
	UsersCollection     = "users"
	ProductsCollection  = "products"
	PurchasesCollection = "purchases"
)

func NewMongoDB(cfg *config.Config) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURL))
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.DatabaseName)

	redis, err := NewRedisDB(cfg)
	if err != nil {
		return nil, err
	}

	mongo := &MongoDB{
		client:   client,
		database: db,
		config:   cfg,
		redis:    redis,
	}

	if err := mongo.createIndexes(ctx); err != nil {
		return nil, err
	}

	return mongo, nil
}

func (m *MongoDB) createIndexes(ctx context.Context) error {
	// Users indexes
	_, err := m.database.Collection(UsersCollection).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})

	return err
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoDB) GetJWTSecret() string {
	return m.config.JWTSecret
}

func (m *MongoDB) Close(ctx context.Context) error {
	if err := m.client.Disconnect(ctx); err != nil {
		return err
	}
	return m.redis.Close()
}

// Collection helpers
func (m *MongoDB) Users() *mongo.Collection {
	return m.database.Collection(UsersCollection)
}

func (m *MongoDB) Products() *mongo.Collection {
	return m.database.Collection(ProductsCollection)
}

func (m *MongoDB) Purchases() *mongo.Collection {
	return m.database.Collection(PurchasesCollection)
}
