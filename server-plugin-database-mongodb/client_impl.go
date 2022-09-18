package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"twowls.org/patchwork/commons/database"
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
	conn   connstring.ConnString
	client *mongo.Client
	db     *mongo.Database
	log    logging.Facade
}

// extension.Extension methods

func (ext *ClientExtension) Configure(opts *extension.Options) error {
	ext.log = opts.Log

	uri := opts.StrConfigDefault("uri", "")
	if uri == "" {
		ext.log.Error().Msg("Connection uri is empty")
		return ErrInvalidUri
	}

	conn, err := connstring.ParseAndValidate(uri)
	if err != nil {
		ext.log.Error().Err(err).Msg("Invalid connection URI")
		return ErrInvalidUri
	} else if conn.Database == "" {
		ext.log.Error().Msg("Connection URI does not include database name")
		return ErrInvalidUri
	}

	clientOpts := options.Client().ApplyURI(uri)
	if username := opts.StrConfigDefault("username", ""); username != "" {
		password, hasPassword := opts.StrConfig("password")
		clientOpts.SetAuth(options.Credential{
			Username:    username,
			Password:    password,
			PasswordSet: hasPassword,
		})
	}

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		ext.log.Error().Err(err).Msg("Create client failed")
		return ErrCreateClient
	}

	ext.db = client.Database(conn.Database, options.Database())
	ext.client = client
	ext.conn = conn
	return nil
}

// database.Client methods

func (ext *ClientExtension) Connect(ctx context.Context) error {
	if err := ext.client.Connect(ctx); err != nil {
		ext.log.Error().Err(err).Msg("Connection error")
		return ErrConnect
	}

	var info bson.M
	if err := ext.db.RunCommand(ctx, bson.D{{"buildInfo", 1}}).Decode(&info); err != nil {
		ext.log.Error().Err(err).Msg("Cannot get server info")
		return ErrConnect
	}

	ext.log.Info().Msgf("Connected to deployment: %v, version %v", ext.conn.Hosts, info["version"])
	return nil
}

func (ext *ClientExtension) Disconnect(ctx context.Context) error {
	if err := ext.client.Disconnect(ctx); err != nil {
		ext.log.Warn().Err(err).Msg("Disconnect error")
		return ErrDisconnect
	}
	return nil
}

func (ext *ClientExtension) CallInTransaction(ctx context.Context, worker database.TxWorkerCallable) (any, error) {
	if worker == nil {
		ext.log.Warn().Msg("CallInTransaction() no worker")
		return nil, nil
	}

	if session, err := ext.client.StartSession(); err == nil {
		var txErr error

		defer func() {
			if txErr != nil {
				if abortErr := session.AbortTransaction(ctx); abortErr != nil {
					ext.log.Error().Err(abortErr).Msg("CallInTransaction() abort transaction failed")
				} else {
					ext.log.Debug().Msg("CallInTransaction() transaction aborted")
				}
			}
			session.EndSession(ctx)
		}()

		sc := mongo.NewSessionContext(ctx, session)
		if err = sc.StartTransaction(ext.txOptions()); err != nil {
			ext.log.Error().Err(err).Msg("CallInTransaction() start transaction failed")
			return nil, database.ErrClientTxFail
		}

		ext.log.Debug().Msg("CallInTransaction() begin transaction")
		var result any
		if result, txErr = worker(sc); txErr != nil {
			ext.log.Error().Err(txErr).Msg("CallInTransaction() worker returned error")
			return nil, txErr
		}

		if txErr = sc.CommitTransaction(sc); txErr != nil {
			ext.log.Error().Err(txErr).Msg("CallInTransaction() commit failed")
			return nil, database.ErrClientTxFail
		}

		ext.log.Debug().Msg("CallInTransaction() transaction committed")
		return result, nil
	} else {
		ext.log.Error().Err(err).Msg("CallInTransaction() start session failed")
		return nil, database.ErrClientTxFail
	}
}

func (ext *ClientExtension) RunInTransaction(ctx context.Context, worker database.TxWorker) error {
	_, err := ext.CallInTransaction(ctx, func(tc context.Context) (any, error) {
		return nil, worker(tc)
	})
	return err
}

// private

func (ext *ClientExtension) txOptions() *options.TransactionOptions {
	return options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
}
