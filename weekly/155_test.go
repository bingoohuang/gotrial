package weekly

import (
	"context"
	"testing"
)

// https://github.com/polaris1119/golangweekly/blob/master/docs/issue-155.md

func Test155(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 0)
	<-ctx.Done()
	t.Log("timed out")
}
