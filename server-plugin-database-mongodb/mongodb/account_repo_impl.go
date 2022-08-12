package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"twowls.org/patchwork/commons/database/repos"
	"twowls.org/patchwork/commons/service"
)

const userAccountCollectionName = "account.user"

// database.repos.AccountRepository methods

func (ext *ClientExtension) AccountFindUser(login string, lookupByEmail bool) *service.AccountUser {
	coll := ext.userAccountCollection()

	var filter bson.D
	if lookupByEmail {
		filter = bson.D{
			{"$or", bson.A{
				bson.M{"loginOrEmail": login},
				bson.M{"email": login}},
			},
		}
	} else {
		filter = bson.D{{"login", login}}
	}

	var account service.AccountUser
	if err := coll.FindOne(context.TODO(), filter).Decode(&account); err != nil {
		if err != mongo.ErrNoDocuments {
			ext.log.Error("AccountFindUser(): query failed: %v", err)
		}
		return nil
	}

	return &account
}

func (ext *ClientExtension) AccountFindLoginUser(loginOrEmail string, passwordMatcher repos.PasswordMatcher) (*service.AccountUser, bool) {
	coll := ext.userAccountCollection()

	filter := bson.D{
		{"$or", bson.A{
			bson.M{"login": loginOrEmail},
			bson.M{"email": loginOrEmail}},
		},
		{"flags", bson.M{
			"$nin": bson.A{
				service.AccountUserInternal,  // excluding system accounts
				service.AccountUserSuspended, // excluding suspended accounts
			},
		}},
	}

	var rawDoc bson.M
	rawResult := coll.FindOne(context.TODO(), filter)
	if err := rawResult.Decode(&rawDoc); err == nil {
		if pwd, ok := rawDoc["pwd"].(primitive.Binary); ok {
			var account service.AccountUser
			if err = rawResult.Decode(&account); err != nil {
				ext.log.Error("AccountFindLoginUser(): failed to decode AccountUser data: %v", err)
			}
			return &account, passwordMatcher != nil && passwordMatcher(pwd.Data)
		}
	} else {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			ext.log.Error("AccountFindLoginUser(): query failed: %v", err)
		}
	}

	return nil, false
}

// private

func (ext *ClientExtension) userAccountCollection() *mongo.Collection {
	coll := ext.db.Collection(userAccountCollectionName)
	indices := []mongo.IndexModel{
		{
			Keys:    bson.M{"login": 1},
			Options: options.Index().SetUnique(true),
		}, {
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	if _, err := coll.Indexes().CreateMany(context.TODO(), indices); err != nil {
		ext.log.Error("userAccountCollection() could not create indices on %q: %v",
			userAccountCollectionName, err)
	}

	return coll
}
