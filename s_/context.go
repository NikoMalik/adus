package s_

import "context"

type AdusKey int

const (
	AdusKeyContext AdusKey = iota
)

func FromContext(ctx context.Context) *Instance {
	if s, ok := ctx.Value(AdusKeyContext).(*Instance); ok {
		return s
	}

	return nil
}
