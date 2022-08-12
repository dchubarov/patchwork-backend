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

func (ext *ClientExtension) AuthFindSession(sid string) *service.AuthSession {
	oid, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		ext.log.Error("AuthFindSession() could not convert sid %q to ObjectID: %v", sid, err)
		return nil
	}

	filter := bson.D{
		{"_id", oid},
		{"expires", bson.M{"$gt": time.Now().UTC()}},
	}

	var sessionBson bson.M
	sessionResult := ext.sessionCollection().FindOne(context.TODO(), filter)
	if err = sessionResult.Decode(&sessionBson); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			ext.log.Error("AuthFindSession() query failed: %v", err)
		}
		return nil
	}

	var session service.AuthSession
	if err = sessionResult.Decode(&session); err != nil {
		ext.log.Error("AuthFindSession() could not decode session result: %v", err)
		return nil
	}

	session.Sid = sid
	return &session
}

func (ext *ClientExtension) AuthNewSession(user *service.AccountUser) *service.AuthSession {
	if user == nil || user.IsSuspended() {
		ext.log.Error("AuthNewSession(): user %q is suspended", user.Login)
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

	if result, err := ext.sessionCollection().InsertOne(context.TODO(), sessionBson); err != nil || result.InsertedID == nil {
		ext.log.Error("AuthNewSession(): insert failed: %v", err)
		return nil
	} else {
		session.Sid = result.InsertedID.(primitive.ObjectID).Hex()
		return session
	}
}

func (ext *ClientExtension) AuthDeleteSession(session *service.AuthSession) bool {
	oid, err := primitive.ObjectIDFromHex(session.Sid)
	if err != nil {
		ext.log.Error("AuthDeleteSession() could not convert sid %q to ObjectID: %v", session.Sid, err)
		return false
	}

	filter := bson.D{{"_id", oid}}
	if result, err := ext.sessionCollection().DeleteOne(context.TODO(), filter); err != nil {
		ext.log.Error("AuthDeleteSession() query failed: %v", err)
		return false
	} else {
		return result.DeletedCount == 1
	}
}

// private

func (ext *ClientExtension) sessionCollection() *mongo.Collection {
	coll := ext.db.Collection(sessionCollectionName)
	index := mongo.IndexModel{Keys: bson.M{"expires": 1}}

	if _, err := coll.Indexes().CreateOne(context.TODO(), index); err != nil {
		ext.log.Error("sessionCollection() could not create index on %q: %v", sessionCollectionName, err)
	}

	return coll
}
