package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"im/Initialize"
	"im/global"
	"im/handle"
	"im/middlewares"
	"log"
	"net/http"
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
	r.Use(middlewares.Cors())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"msg":"ok",
			"code":"200",
		})
	})
	r.POST("/api/register",handle.Register)
	r.GET("/api/list",handle.UserList)
	r.POST("/api/login",handle.Login)
	r.GET("/api/ws",handle.WS)
	r.GET("/api/loginlist",handle.LoginList)
	r.Use(middlewares.JWT())
	r.GET("/api/user",handle.UserInfo)

	// 启动运行
	if err := r.Run(fmt.Sprintf(":%d", global.SrvConfig.Port)); err != nil {
		log.Fatal("启动失败：",err)
	}
}
