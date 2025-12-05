package main

import (
	"ginchat/config"
	_ "ginchat/docs"
	"ginchat/models"
	"ginchat/router"
	"ginchat/utils"
	"github.com/spf13/viper"
	"time"
)

// @title GinChat
// @version 1.0

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	InitTimer()

	r := router.Router()
	r.Run(viper.GetString("port.server"))
}

// 初始化定时器（使用全局配置）
func InitTimer() {
	// 从全局配置获取超时参数
	delay := time.Duration(config.GlobalConfig.Timeout.DelayHeartbeat) * time.Second
	tick := time.Duration(config.GlobalConfig.Timeout.HeartbeatHz) * time.Second
	utils.Timer(delay, tick, models.CleanConnection, "")
}
