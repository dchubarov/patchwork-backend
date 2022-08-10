package database

import "context"

// Client represents a database client
type Client interface {
	// Connect establishes connection to the database
	Connect(context.Context) error

	// Disconnect frees resources and shutdowns database connection
	Disconnect(context.Context) error
}
