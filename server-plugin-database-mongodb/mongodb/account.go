package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"twowls.org/patchwork/commons/database/domain"
)

// database.repos.AccountRepository methods

func (ext *ClientExtension) AccountFindUser(login string) domain.UserAccount {
	coll := ext.db.Collection("account.users", options.Collection())

	var result domain.UserAccount
	if err := coll.FindOne(context.TODO(), bson.D{{"login", login}}).Decode(&result); err != nil {
		if err != mongo.ErrNoDocuments {
			ext.log.Error("AccountFindUser(): query failed: %v", err)
		}
		return domain.EmptyUserAccount
	}

	result.Exists = true
	return result
}
