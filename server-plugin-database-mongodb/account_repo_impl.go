package main

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

func (ext *ClientExtension) AccountFindUser(ctx context.Context, login string, lookupByEmail bool) *service.UserAccount {
	coll := ext.userAccountCollection(ctx)

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

	var account service.UserAccount
	if err := coll.FindOne(ctx, filter).Decode(&account); err != nil {
		if err != mongo.ErrNoDocuments {
			ext.log.Error().Err(err).Msg("AccountFindUser(): query failed")
		}
		return nil
	}

	return &account
}

func (ext *ClientExtension) AccountFindLoginUser(ctx context.Context, loginOrEmail string, passwordMatcher repos.PasswordMatcher) (*service.UserAccount, bool) {
	coll := ext.userAccountCollection(ctx)

	filter := bson.D{
		{"$or", bson.A{
			bson.M{"login": loginOrEmail},
			bson.M{"email": loginOrEmail}},
		},
		{"flags", bson.M{
			"$nin": bson.A{
				service.UserAccountInternal,  // excluding system accounts
				service.UserAccountSuspended, // excluding suspended accounts
			},
		}},
	}

	var rawDoc bson.M
	rawResult := coll.FindOne(ctx, filter)
	if err := rawResult.Decode(&rawDoc); err == nil {
		if pwd, ok := rawDoc["pwd"].(primitive.Binary); ok {
			var account service.UserAccount
			if err = rawResult.Decode(&account); err != nil {
				ext.log.Error().Err(err).Msg("AccountFindLoginUser(): failed to decode UserAccount data")
			}
			return &account, passwordMatcher != nil && passwordMatcher(pwd.Data)
		}
	} else {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			ext.log.Error().Err(err).Msg("AccountFindLoginUser(): query failed")
		}
	}

	return nil, false
}

// private

func (ext *ClientExtension) userAccountCollection(ctx context.Context) *mongo.Collection {
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

	if _, err := coll.Indexes().CreateMany(ctx, indices); err != nil {
		ext.log.Error().Err(err).Msgf("userAccountCollection() could not create indices on %q",
			userAccountCollectionName)
	}

	return coll
}
