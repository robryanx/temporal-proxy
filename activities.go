package synchronousproxy

import (
	"context"

	"go.temporal.io/sdk/activity"
)

func RegisterEmail(ctx context.Context, email string) error {
	logger := activity.GetLogger(ctx)

	logger.Info("activity: registered email", email)

	return nil
}
