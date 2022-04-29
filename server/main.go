package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"im/Initialize"
	"im/global"
	"im/handle"
	"im/middlewares"
	"log"
)

func main() {
	Initialize.InitConfig()
	Initialize.InitUserMap()
	Initialize.InitRedis()
	global.UserMap.Range(func(key, value interface{}) bool {
		fmt.Println(key)
		fmt.Println(value)
		//fmt.Println(len(key.(string)))
		return true
	})
	r := gin.Default()
	// 跨域
	r.Use(middlewares.Cors())

	// 注册
	r.POST("/api/register", handle.Register)
	// 已注册用户列表
	r.GET("/api/list", handle.UserList)
	// 登录
	r.POST("/api/login", handle.Login)
	// ws连接
	r.GET("/api/ws", handle.WS)
	// 获取登录列表(目前没用到)
	r.GET("/api/loginlist", handle.LoginList)
	// JWT
	r.Use(middlewares.JWT())
	// 获取用户名
	r.GET("/api/user", handle.UserInfo)

	// 启动运行
	if err := r.Run(fmt.Sprintf(":%d", global.SrvConfig.Port)); err != nil {
		log.Fatal("启动失败：", err)
	}
}
