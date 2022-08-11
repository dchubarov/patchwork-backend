package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"twowls.org/patchwork/commons/database/repos"
)

const sessionCollectionName = "auth.session"
const sessionTimeToLive = time.Hour * 8 // TODO session ttl must be configurable

var (
	ErrAuthUserCantOwnSession  = errors.New("invalid or suspended user")
	ErrAuthCannotCreateSession = errors.New("could not create session")
	ErrAuthSessionNotFound     = errors.New("session not found")
)

// database.repos.AuthRepository methods

func (ext *ClientExtension) AuthFindSession(sid string) (*repos.AuthSession, error) {
	if oid, err := primitive.ObjectIDFromHex(sid); err == nil {
		filter := bson.D{
			{"_id", oid},
			{"expire_utc", bson.M{"$lt": time.Now().UTC()}},
		}

		var sessionBson bson.M
		if err := ext.sessionCollection().FindOne(context.TODO(), filter).Decode(&sessionBson); err == nil {
			return &repos.AuthSession{
				Created: sessionBson["create_utc"].(time.Time),
				Expires: sessionBson["expire_utc"].(time.Time),
				Sid:     sid,
			}, nil
		} else {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				ext.log.Error("AuthFindSession() query failed: %v", err)
			}
			return nil, ErrAuthSessionNotFound
		}
	} else {
		ext.log.Error("AuthFindSession() could not convert %q to ObjectID", sid)
		return nil, ErrAuthSessionNotFound
	}
}

func (ext *ClientExtension) AuthNewSession(user *repos.AccountUser) (*repos.AuthSession, error) {
	if user == nil || user.IsSuspended() {
		return nil, ErrAuthUserCantOwnSession
	}

	timestamp := time.Now().UTC()
	session := &repos.AuthSession{
		Created: timestamp,
		Expires: timestamp.Add(sessionTimeToLive),
	}

	sessionBson := bson.D{
		{"user", user.Login},
		{"privileged", user.IsPrivileged()},
		{"create_utc", session.Created},
		{"expire_utc", session.Expires},
	}

	if result, err := ext.sessionCollection().InsertOne(context.TODO(), sessionBson); err != nil || result.InsertedID == nil {
		ext.log.Error("AuthNewSession(): insert failed: %v", err)
		return nil, ErrAuthCannotCreateSession
	} else {
		session.Sid = result.InsertedID.(primitive.ObjectID).Hex()
		return session, nil
	}
}

func (ext *ClientExtension) sessionCollection() *mongo.Collection {
	return ext.db.Collection(sessionCollectionName)
}
