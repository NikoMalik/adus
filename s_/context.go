package s_

import "context"

type AdusKey int

const (
	AdusKeyContext AdusKey = iota
)

func FromContext(ctx context.Context)
