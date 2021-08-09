package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// RPush : 資料寫入 List 表尾
func RPush(ctx context.Context, key string, value interface{}) error {

	return GetConn("Redis").RPush(
		ctx,
		key,
		value,
	).Err()
}

// LPos : 找資料在 List 的 index
func LPos(ctx context.Context, key string, value string) (int64, error) {

	return GetConn("Redis").LPos(
		ctx,
		key,
		value,
		redis.LPosArgs{
			Rank:   0,
			MaxLen: 0,
		},
	).Result()
}

// LLen : 返回列表 key 的長度
func LLen(ctx context.Context, key string) (int64, error) {

	return GetConn("Redis").LLen(
		ctx,
		key,
	).Result()
}

// LRange : 返回列表 key 中指定區間內的元素，區間以 start 和 stop 指定 (全部 start 0 stop -1)
func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {

	return GetConn("Redis").LRange(
		ctx,
		key,
		start,
		stop,
	).Result()
}

// LRem : 從列表 key 中刪除前 count 個數等於 value 的元素，count = 0 移除所有值為 value 的元素
func LRem(ctx context.Context, key string, count int64, value interface{}) error {

	return GetConn("Redis").LRem(
		ctx,
		key,
		count,
		value,
	).Err()
}
