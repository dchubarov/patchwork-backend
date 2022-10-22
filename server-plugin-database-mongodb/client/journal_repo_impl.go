package client

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const journalEventCollectionName = "journal.event"

var emptyDetails = bson.M{}

func (ext *ClientExtension) JournalAddEvent(ctx context.Context, event string, user string, details map[string]any) {
	if len(details) < 1 {
		details = emptyDetails
	}

	document := bson.D{
		{"timestamp", time.Now()},
		{"event", event},
		{"user", user},
		{"details", details},
	}

	_, err := ext.db.Collection(journalEventCollectionName).InsertOne(ctx, document)
	if err != nil {
		ext.log.Error().Err(err).Msg("JournalAddEvent() insert failed")
	}
}
