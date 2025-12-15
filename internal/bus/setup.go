package bus

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

func Redis_Init(rdb *redis.Client) error {
	ctx := context.Background()
	err := rdb.XGroupCreateMkStream(ctx, "twitch:events", "workers", "$").Err()

	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}
	return nil
}
