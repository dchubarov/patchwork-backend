package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"twowls.org/patchwork/backend/config"
	"twowls.org/patchwork/backend/logging"
)

type Client struct {
	mongo *mongo.Client
}

func (c *Client) Connect(ctx context.Context) {
	if err := c.mongo.Connect(ctx); err != nil {
		logging.Panic("[mongo] cannot establish connection to mongo")
	}

	if err := c.mongo.Ping(ctx, readpref.Primary()); err != nil {
		logging.Panic("[mongo] ping failed")
	}

	logging.Info("[mongo] connected")
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.mongo.Disconnect(ctx)
}

func New(cfg config.Database) *Client {
	opts := options.Client().ApplyURI(cfg.Url)
	if cfg.Username != "" {
		opts.SetAuth(options.Credential{
			Username:    cfg.Username,
			Password:    cfg.Password,
			PasswordSet: cfg.Password != "",
		})
	}

	c, err := mongo.NewClient(opts)
	if err != nil {
		logging.Panic("mongo.NewClient() failed: %v", err)
	}

	return &Client{mongo: c}
}
