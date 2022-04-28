package global

import (
	"github.com/go-redis/redis/v8"
	"im/config"
	"sync"
)

var (
	//配置信息
	SrvConfig = config.SrvConfig{}
	//已注册用户map，key为name value为password
	UserMap = sync.Map{}
	//在线用户map，key为name value为连接句柄list
	LoginMap = sync.Map{}
	//redis客户端
	Redis *redis.Client
)


