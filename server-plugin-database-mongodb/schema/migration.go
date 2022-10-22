package schema

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterMigrationFunc adds a migration function
func RegisterMigrationFunc(f func(context.Context, *mongo.Database) error) {

}

// RegisterMigrationCommand adds a migration command
func RegisterMigrationCommand(cmd bson.D) {
	RegisterMigrationFunc(func(ctx context.Context, db *mongo.Database) error {
		result := db.RunCommand(ctx, cmd)

		if err := result.Err(); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}

		return nil
	})
}
