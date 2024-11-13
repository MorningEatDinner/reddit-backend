package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/xiaorui/reddit-async/reddit-backend/settings"
	"go.uber.org/zap"
)

type RedisClient struct {
	Client  *redis.Client
	Context context.Context
}

// var rdb *redis.Client

var RDB *RedisClient

func Init(cfg *settings.RedisConfig) (err error) {
	// rdb = redis.NewClient(&redis.Options{
	// 	Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	// 	Password: cfg.Password,
	// 	DB:       cfg.DB,
	// 	PoolSize: cfg.PoolSize,
	// })
	// ctx := context.Background()
	// _, err = rdb.Ping(ctx).Result()
	// return
	RDB = &RedisClient{}
	RDB.Context = context.Background()
	RDB.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	_, err = RDB.Client.Ping(RDB.Context).Result()
	zap.L().Info("测试连接", zap.Int("链接成功", cfg.Port))
	return
}

func Close() {
	_ = RDB.Client.Close()
}
