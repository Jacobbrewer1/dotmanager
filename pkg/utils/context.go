package utils

import "context"

func CoreContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
