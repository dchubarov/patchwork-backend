package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const journalEventCollectionName = "journal.event"

func (ext *ClientExtension) JournalAddEvent(ctx context.Context, event string, user string, data map[string]any) {
	document := bson.D{
		{"event", event},
		{"user", user},
		{"created", time.Now().UTC()},
	}

	if len(data) > 0 {
		document = append(document, bson.E{Key: "data", Value: data})
	}

	_, err := ext.db.Collection(journalEventCollectionName).InsertOne(ctx, document)
	if err != nil {
		ext.log.Error().Err(err).Msg("JournalAddEvent() insert failed")
	}
}
