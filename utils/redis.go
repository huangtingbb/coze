package utils

import (
	"context"
	"coze-agent-platform/config"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() {
	redisCfg := config.Cfg.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	// 测试Redis连接
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败: %v", err)
		log.Println("警告: Redis不可用，某些功能可能无法正常工作")
		return
	}

	log.Println("Redis连接初始化完成")
}

// GetRedisClient 获取Redis客户端实例
func GetRedisClient() *redis.Client {
	return RDB
}
