package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	//"ginchat/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited.")

}

func InitMySQL() {
	//自定义日志模板，打印SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢SQL阈值
			LogLevel:      logger.Info, //级别
			Colorful:      true,        //彩色
		},
	)

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{Logger: newLogger})
	/*if err != nil {
		panic("failed to connect database")
	}*/
	//user := models.UserBasic{}
	//DB.Find(&user)
	fmt.Println("config mysql inited.")
}

func InitRedis() {

	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	pong, err := Red.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("config redis init unsuccessfully: ", err)
	} else {
		fmt.Println("config redis inited: ", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	err = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println("redis publish unsuccessfully: ", err)
	}
	return err
}

// Subscribe 订阅redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Red.Subscribe: ", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println("Red.Subscribe unsuccessfully: ", err)
		return "", err
	}
	fmt.Println("sub receive msg: ", msg.Payload)
	return msg.Payload, err
}
