package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"twowls.org/patchwork/plugin/database/mongodb/schema"
)

var _ = bson.D{
	{"create", "account.user"},
	{"comment", "User accounts"},
	{"validator", bson.D{}},
}

func createAccountUserCollection(ctx context.Context, db *mongo.Database) error {
	err := db.CreateCollection(ctx, schema.UserAccountCollection, options.CreateCollection().
		SetValidator(""))

	if err != nil {
		return err
	}

	db.RunCommand(ctx, bson.D{
		{"createIndexes", schema.UserAccountCollection},
		{"indexes", bson.A{
			bson.D{
				{"key", bson.D{{"login", 1}}},
				{"unique", true},
			},
		}},
	})

	return nil
}

func init() {
	schema.RegisterMigrationCommand(bson.D{
		{"create", schema.UserAccountCollection},
		{"comment", "User accounts"},
	})
}
