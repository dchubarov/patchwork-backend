package database

import (
	"context"
	"errors"
)

var (
	ErrClientTxFail = errors.New("transaction failure")
)

type TxWorkerCallable func(context.Context) (any, error)
type TxWorker func(context.Context) error

// Client represents a database client
type Client interface {
	// Connect establishes connection to the database
	Connect(context.Context) error
	// Disconnect frees resources and shutdowns database connection
	Disconnect(context.Context) error
	// CallInTransaction executes worker function in transaction and returns worker result
	CallInTransaction(ctx context.Context, worker TxWorkerCallable) (any, error)
	// RunInTransaction executes worker within a transaction
	RunInTransaction(ctx context.Context, worker TxWorker) error
}
