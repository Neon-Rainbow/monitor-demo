package config

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	conf = &Config{}
	once sync.Once
)

type Config struct {
	Server     *Server
	Logger     *Logger
	Prometheus *Prometheus
	Etcd       *Etcd
	Consul     *Consul
	Grafana    *Grafana
}

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type Logger struct {
	Level      int    `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputPath string `mapstructure:"output_path"`
}

type Prometheus struct {
	Enable bool
}

type Etcd struct {
	Endpoints []string `mapstructure:"endpoints"`
	Timeout   int      `mapstructure:"timeout"`
}

type Consul struct {
	Enable bool
}

type Grafana struct {
	Enable bool
}

// initConfig 用于初始化配置文件
func initConfig() {
	// 设置配置文件名
	viper.SetConfigFile("config.toml")

	// 设置默认配置
	setDefault()

	// 用于判断配置文件是否被修改
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件被修改:", e.Name)
	})

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %s", err))
	}

	// 将配置文件内容解析到结构体中
	if err := viper.Unmarshal(conf); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %s", err))
	}
}

// Get 用于获取配置文件
// 通过 sync.Once 来保证只初始化一次
//
// 返回值:
//   - *Config: 配置文件结构体指针
func Get() *Config {
	once.Do(func() {
		initConfig()
	})
	return conf
}

func setDefault() {
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.log_level", 0)

	viper.SetDefault("logger.level", 0)
	viper.SetDefault("logger.format", "console")
	viper.SetDefault("logger.output_path", "./logs/")
}
