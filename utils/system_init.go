package utils

import (
	"context"
	"fmt"
	"ginchat/config"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

// 初始化配置（加载到全局结构体）
func InitConfig() {
	// 设置配置文件
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %v", err))
	}

	// 支持环境变量覆盖（例如：将配置文件中的mysql.dsn映射为环境变量MYSQL_DSN）
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ginchat")                          // 环境变量前缀，避免冲突
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 将点分隔转换为下划线

	// 解析配置到全局结构体
	if err := viper.Unmarshal(&config.GlobalConfig); err != nil {
		panic(fmt.Sprintf("解析配置失败: %v", err))
	}

	fmt.Println("配置加载完成:", config.GlobalConfig)
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

	// 使用全局配置中的MySQL DSN
	db, err := gorm.Open(mysql.Open(config.GlobalConfig.MySQL.DSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(fmt.Sprintf("初始化MySQL失败: %v", err))
	}

	// 设置连接池（新增）
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(config.GlobalConfig.MySQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.GlobalConfig.MySQL.MaxIdleConns)

	DB = db
	fmt.Println("MySQL初始化完成")
}

func InitRedis() {

	Red = redis.NewClient(&redis.Options{
		Addr:         config.GlobalConfig.Redis.Addr,
		Password:     config.GlobalConfig.Redis.Password,
		DB:           config.GlobalConfig.Redis.DB,
		PoolSize:     config.GlobalConfig.Redis.PoolSize,
		MinIdleConns: config.GlobalConfig.Redis.MinIdleConns,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	pong, err := Red.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("初始化Redis失败: %v", err))
	}
	fmt.Println("Redis初始化完成:", pong)
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
