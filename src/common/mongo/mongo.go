package appmongo

import (
	"context"
	"fmt"
	"go-kit/src/common/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DBProvider struct {
	db     *mongo.Database
	client *mongo.Client
}

func (p *DBProvider) Start(ctx context.Context) error {
	if err := p.client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("[MongoDB] failed to ping to DB %w", err)
	}
	return nil
}

func (p *DBProvider) Stop(ctx context.Context) error {
	err := p.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("[MongoDB] failed to disconnect %w", err)
	}
	return nil
}

func (p *DBProvider) DB() *mongo.Database {
	return p.db
}

func NewMongoProvider(cf *configs.Config) (*DBProvider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.ClientOptions{}

	uri := cf.Mongo.Uri
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri), &opts)
	if err != nil {
		return nil, fmt.Errorf("[MongoDB] failed to connect to DB %w", err)
	}

	db := client.Database(cf.Mongo.DB)
	return &DBProvider{
		db:     db,
		client: client,
	}, nil
}
