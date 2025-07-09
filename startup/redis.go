package startup

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// InitializeRedis Redis 클라이언트를 초기화하고 연결을 테스트합니다
func InitializeRedis(ctx context.Context, redisAddr string) (*redis.Client, error) {
	log := logrus.New()
	log.SetReportCaller(true)

	// Redis 클라이언트 생성
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// 연결 테스트
	err := rdb.Set(ctx, "hello", "world", 0).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set test key in Redis: %w", err)
	}

	val, err := rdb.Get(ctx, "hello").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get test key from Redis: %w", err)
	}

	log.Infof("Redis connection test successful: hello=%s", val)
	return rdb, nil
}
