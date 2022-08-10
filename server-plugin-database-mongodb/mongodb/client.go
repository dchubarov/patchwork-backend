package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"twowls.org/patchwork/commons/extension"
	"twowls.org/patchwork/commons/logging"
)

var (
	ErrInvalidUri   = errors.New("invalid connection string")
	ErrCreateClient = errors.New("failed to initialize database client")
	ErrConnect      = errors.New("failed to connect to database")
	ErrDisconnect   = errors.New("failed to disconnect from database")
)

type ClientExtension struct {
	client *mongo.Client
	db     *mongo.Database
	log    logging.Facade
}

// extension.Extension methods

func (ext *ClientExtension) Configure(opts *extension.Options) error {
	ext.log = opts.Log

	uri := opts.StrConfigDefault("uri", "")
	if uri == "" {
		ext.log.Error("Connection uri is empty")
		return ErrInvalidUri
	}

	conn, err := connstring.ParseAndValidate(uri)
	if err != nil {
		ext.log.Error("Invalid connection URI: %v", err)
		return ErrInvalidUri
	} else if conn.Database == "" {
		ext.log.Error("Connection URI does not include database name")
		return ErrInvalidUri
	}

	clientOpts := options.Client().ApplyURI(uri)
	if username, ok := opts.StrConfig("username"); ok {
		password, hasPassword := opts.StrConfig("password")
		clientOpts.SetAuth(options.Credential{
			Username:    username,
			Password:    password,
			PasswordSet: hasPassword,
		})
	}

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		ext.log.Error("Create client failed: %v", err)
		return ErrCreateClient
	}

	ext.db = client.Database(conn.Database, options.Database())
	ext.client = client
	return nil
}

// database.Client methods

func (ext *ClientExtension) Connect(ctx context.Context) error {
	if err := ext.client.Connect(ctx); err != nil {
		ext.log.Error("Connection error: %v", err)
		return ErrConnect
	}

	if err := ext.db.Client().Ping(ctx, readpref.Primary()); err != nil {
		ext.log.Error("Unable to ping database server: %v", err)
		return ErrConnect
	}

	ext.log.Info("Connected")
	return nil
}

func (ext *ClientExtension) Disconnect(ctx context.Context) error {
	if err := ext.client.Disconnect(ctx); err != nil {
		ext.log.Error("Disconnect error %v", err)
		return ErrDisconnect
	}
	return nil
}
