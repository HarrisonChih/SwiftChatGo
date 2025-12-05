package main

import (
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

// 初始化定时器
func InitTimer() {
	utils.Timer(time.Duration(viper.GetInt("timeout.DelayHeartbeat"))*time.Second, time.Duration(viper.GetInt("timeout.HeartbeatHz"))*time.Second, models.CleanConnection, "")
}
