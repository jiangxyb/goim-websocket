package Initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"im/global"
)

func InitConfig() {
	//从配置文件中读取出对应的配置
	var configFileName = fmt.Sprintf("./config.yaml" )
	v := viper.New()
	//文件的路径
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 开启实时监控
	v.WatchConfig()
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(&global.SrvConfig); err != nil {
		panic(err)
	}

	// 文件更新的回调函数
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置改变")
		if err := v.Unmarshal(&global.SrvConfig); err != nil {
			panic(err)
		}
	})
}

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.IsSet(env)
}


