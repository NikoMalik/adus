package ctx

import "context"

type SessionKey int

type ID uint64

const (
	SessionKeyID SessionKey = iota
)

// ContextWithId returns a new context with the given id as the value
// associated with the SessionKeyID key.
func ContextWithId(ctx context.Context, id ID) context.Context {
	return context.WithValue(ctx, SessionKeyID, id)
}

// FromContext returns the id associated with the SessionKeyID key from the
// given context. If there is no value associated with the key, it returns 0.
func FromContext(ctx context.Context) ID {
	if id, ok := ctx.Value(SessionKeyID).(ID); ok {
		return id
	}
	return 0
}
