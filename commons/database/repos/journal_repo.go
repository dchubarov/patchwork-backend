package repos

import "context"

// JournalRepository contains methods related to journal
type JournalRepository interface {
	// JournalAddEvent adds a new event to journal
	JournalAddEvent(ctx context.Context, event string, user string, data map[string]any)
}
