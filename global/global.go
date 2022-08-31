package global

import (
	"chatroom/config"

	"github.com/go-redis/redis/v8"
)

var (
	GVA_CONFIG config.Server
	GVA_REDIS  *redis.Client
)
