package config

// 读取yaml配置文件,形成映射的相关类
type JWTconfig struct {
	PrivateKeyPath string `mapstructure:"privateKeyPath" json:"privateKeyPath"`
	PublicKeyPath string `mapstructure:"publicKeyPath" json:"publicKeyPath"`
}

type RedisConfig struct {
	IP string `mapstructure:"ip"`
	Port string `mapstructure:"port"`
	PassWord string `mapstructure:"password"`
}

type SrvConfig struct {
	Name    string    `mapstructure:"name" json:"name"`
	Port    int       `mapstructure:"port" json:"port"`
	JWTInfo JWTconfig `mapstructure:"jwt" json:"jwt"`
	RedisInfo RedisConfig `mapstructure:"redis" json:"redis"`
}







