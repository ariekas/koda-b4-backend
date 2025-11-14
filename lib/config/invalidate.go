package config

import (
	"context"
	"fmt"
)

func InvalidateRedis(path string) {
	redis := Redis()
	iter := redis.Scan(context.Background(), 0, path, 0).Iterator()
	for iter.Next(context.Background()) {
		redis.Del(context.Background(), iter.Val())
	}

	if err := iter.Err(); err != nil {
		fmt.Printf("error: failed to delete cache, %s", err)
	}
}
