package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"twowls.org/patchwork/server/bootstrap/config"
	"twowls.org/patchwork/server/bootstrap/logging"
)

type Client struct {
	db *mongo.Database
}

func (c *Client) Connect(ctx context.Context) {
	if err := c.db.Client().Connect(ctx); err != nil {
		logging.Panic("[mongo] cannot establish connection to mongo")
	}

	if err := c.db.Client().Ping(ctx, readpref.Primary()); err != nil {
		logging.Panic("[mongo] ping failed: %v", err)
	}

	// TODO remove
	err := c.db.CreateCollection(ctx, "___test__", options.CreateCollection())
	if err != nil {
		logging.Error("error create collection: %v", err)
	}
	// TODO end remove

	logging.Info("[mongo] connected")
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.db.Client().Disconnect(ctx)
}

func New(cfg config.Database) *Client {
	conn, err := connstring.ParseAndValidate(cfg.Url)
	if err != nil {
		logging.Panic("invalid connection string")
	}

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

	return &Client{db: c.Database(conn.Database)}
}
