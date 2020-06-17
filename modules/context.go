package modules

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/micro/go-micro/v2/config"
)

type moduleWrapper struct {
	Redis *redis.Client
}

//ModuleContext 模块上下文
var ModuleContext moduleWrapper

//Setup 初始化Modules
func Setup() {
	//  redis
	var rConfig RedisConfig
	if err := config.Get("config", "redis").Scan(&rConfig); err != nil {
		panic(err)
	}
	ModuleContext.Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", rConfig.Address, rConfig.Port),
		Password: "", // no password set
		DB:       0,
	})
	_, err := ModuleContext.Redis.Ping().Result()
	if err != nil {
		panic(err)
	}
}
