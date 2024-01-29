package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	CODEPLAY_INSTANCE = "codeplay:instance_heartbeat" // sorted set, score为最近心跳时间
)

var ShareCache redis.UniversalClient

func InitShareCache(addr, password string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return err
	}
	ShareCache = rdb
	return nil
}
