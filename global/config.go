/**
 * @Author: cyj19
 * @Date: 2021/12/6 10:37
 */

// 读取配置文件

package global

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	SensitiveWords  []string
	MessageQueueLen = 1024
	TokenSecret     string
)

func initConfig() {
	viper.SetConfigName("chatroom")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(RootDir + "/config")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

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
