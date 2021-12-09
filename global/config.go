/**
 * @Author: cyj19
 * @Date: 2021/12/6 10:37
 */

// 读取配置文件

package global

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"testing"
)

var (
	Addr            string
	SensitiveWords  []string
	MessageQueueLen = 1024
	TokenSecret     string
)

func initConfig() {
	// 增加命令行参数
	var configPath string
	flag.StringVar(&configPath, "config", "", "配置文件所在目录")
	testing.Init()
	flag.Parse()

	viper.SetConfigName("chatroom")
	viper.SetConfigType("yaml")

	if configPath == "" {
		// 默认路径
		configPath = RootDir + "/config"
	}

	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	Addr = viper.GetString("addr")
	SensitiveWords = viper.GetStringSlice("sensitive")
	MessageQueueLen = viper.GetInt("message-queue-len")
	TokenSecret = viper.GetString("token-secret")

	// 热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		_ = viper.ReadInConfig()
		SensitiveWords = viper.GetStringSlice("sensitive")
		TokenSecret = viper.GetString("token-secret")
	})
}
