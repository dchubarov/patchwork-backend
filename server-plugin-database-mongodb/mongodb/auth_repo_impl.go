package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"twowls.org/patchwork/commons/service"
)

const sessionCollectionName = "auth.session"
const sessionTimeToLive = time.Hour * 8 // TODO session ttl must be configurable

// database.repos.AuthRepository methods

func (ext *ClientExtension) AuthFindSession(ctx context.Context, sid string) *service.AuthSession {
	oid, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		ext.log.Error().Err(err).Msgf("AuthFindSession() could not convert sid %q to ObjectID", sid)
		return nil
	}

	filter := bson.D{
		{"_id", oid},
		{"expires", bson.M{"$gt": time.Now().UTC()}},
	}

	var sessionBson bson.M
	sessionResult := ext.sessionCollection(ctx).FindOne(ctx, filter)
	if err = sessionResult.Decode(&sessionBson); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			ext.log.Error().Err(err).Msg("AuthFindSession() query failed")
		}
		return nil
	}

	var session service.AuthSession
	if err = sessionResult.Decode(&session); err != nil {
		ext.log.Error().Err(err).Msg("AuthFindSession() could not decode session result")
		return nil
	}

	session.Sid = sid
	return &session
}

func (ext *ClientExtension) AuthNewSession(ctx context.Context, user *service.UserAccount) *service.AuthSession {
	if user == nil || user.IsSuspended() {
		ext.log.Error().Msgf("AuthNewSession(): user %q is suspended", user.Login)
		return nil
	}

	timestamp := time.Now().UTC()
	session := &service.AuthSession{
		Created: timestamp,
		Expires: timestamp.Add(sessionTimeToLive),
	}

	sessionBson := bson.D{
		{"user", user.Login},
		{"privileged", user.IsPrivileged()},
		{"created", session.Created},
		{"expires", session.Expires},
	}

	if result, err := ext.sessionCollection(ctx).InsertOne(ctx, sessionBson); err != nil || result.InsertedID == nil {
		ext.log.Error().Err(err).Msg("AuthNewSession(): insert failed")
		return nil
	} else {
		session.Sid = result.InsertedID.(primitive.ObjectID).Hex()
		return session
	}
}

func (ext *ClientExtension) AuthDeleteSession(ctx context.Context, session *service.AuthSession) bool {
	oid, err := primitive.ObjectIDFromHex(session.Sid)
	if err != nil {
		ext.log.Error().Err(err).Msgf("AuthDeleteSession() could not convert sid %q to ObjectID", session.Sid)
		return false
	}

	filter := bson.D{{"_id", oid}}
	if result, err := ext.sessionCollection(ctx).DeleteOne(ctx, filter); err != nil {
		ext.log.Error().Err(err).Msg("AuthDeleteSession() delete failed")
		return false
	} else {
		return result.DeletedCount == 1
	}
}

// private

func (ext *ClientExtension) sessionCollection(ctx context.Context) *mongo.Collection {
	coll := ext.db.Collection(sessionCollectionName)
	index := mongo.IndexModel{Keys: bson.M{"expires": 1}}

	if _, err := coll.Indexes().CreateOne(ctx, index); err != nil {
		ext.log.Error().Err(err).Msgf("sessionCollection() could not create index on %q", sessionCollectionName)
	}

	return coll
}
