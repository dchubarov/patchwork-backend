package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"twowls.org/patchwork/commons/database/repos"
)

// database.repos.AccountUserRepository methods

func (ext *ClientExtension) AccountUserFind(loginOrEmail string) (*repos.AccountUser, bool) {
	coll := ext.db.Collection("account.users", options.Collection())

	filter := bson.D{
		{"$or", bson.A{
			bson.M{"loginOrEmail": loginOrEmail},
			bson.M{"email": loginOrEmail}},
		},
	}

	var account repos.AccountUser
	if err := coll.FindOne(context.TODO(), filter).Decode(&account); err != nil {
		if err != mongo.ErrNoDocuments {
			ext.log.Error("AccountUserFind(): query failed: %v", err)
		}
		return nil, false
	}

	return &account, true
}

func (ext *ClientExtension) AccountFindLoginUser(loginOrEmail string, passwordMatcher repos.PasswordMatcher) (*repos.AccountUser, bool) {
	coll := ext.userAccountCollection()

	filter := bson.D{
		{"$or", bson.A{
			bson.M{"login": loginOrEmail},
			bson.M{"email": loginOrEmail}},
		},
		{"flags", bson.M{
			"$nin": bson.A{
				repos.AccountUserInternal,  // excluding system accounts
				repos.AccountUserSuspended, // excluding suspended accounts
			},
		}},
	}

	var rawDoc bson.M
	rawResult := coll.FindOne(context.TODO(), filter)
	if err := rawResult.Decode(&rawDoc); err == nil {
		if pwd, ok := rawDoc["pwd"].(primitive.Binary); ok && passwordMatcher(pwd.Data) {
			var account repos.AccountUser
			if err := rawResult.Decode(&account); err != nil {
				ext.log.Error("AccountFindLoginUser(): failed to decode AccountUser data: %v", err)
			}
			return &account, true
		} else {
			ext.log.Warn("AccountFindLoginUser(): passwords do not match")
		}
	} else {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			ext.log.Error("AccountFindLoginUser(): query failed: %v", err)
		}
	}

	return nil, false
}

func (ext *ClientExtension) userAccountCollection() *mongo.Collection {
	return ext.db.Collection("account.users", options.Collection())
}
