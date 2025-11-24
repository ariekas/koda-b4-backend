package config

import (
	"context"
	"fmt"
)

func InvalidateRedis(pattern string) {
    redis := Redis()
    ctx := context.Background()
    
    var cursor uint64
    var keys []string
    
    for {
        var scanKeys []string
        var err error
        
        scanKeys, cursor, err = redis.Scan(ctx, cursor, pattern, 100).Result()
        if err != nil {
            fmt.Printf("error: failed to scan cache, %s\n", err)
            return
        }
        
        keys = append(keys, scanKeys...)
        
        if cursor == 0 {
            break
        }
    }
    
    if len(keys) > 0 {
        deleted, err := redis.Del(ctx, keys...).Result()
        if err != nil {
            fmt.Printf("error: failed to delete cache, %s\n", err)
        } else {
            fmt.Printf("Invalidated %d cache entries\n", deleted)
        }
    } else {
        fmt.Println("No cache keys found to invalidate")
    }
}