package version

import (
	"context"
	"fmt"
	"runtime"
)

// Set at linking time
var (
	Commit string
	Date   string
)

func Print(ctx context.Context) error {
	fmt.Printf(
		"Commit: %s\nRuntime: %s %s/%s\nDate: %s\n",
		Commit,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		Date,
	)
	return nil
}
