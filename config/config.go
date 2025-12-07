package config

import (
	"time"
)

// 全局配置结构体
type Config struct {
	MySQL   MySQLConfig   `mapstructure:"mysql"`
	Redis   RedisConfig   `mapstructure:"redis"`
	OSS     OSSConfig     `mapstructure:"oss"`
	Timeout TimeoutConfig `mapstructure:"timeout"`
	Port    PortConfig    `mapstructure:"port"`
	App     AppConfig     `mapstructure:"app"` // 新增应用级配置
	JWT     JWTConfig     `mapstructure:"jwt"` // 预留JWT配置
	UDP     UDPConfig     `mapstructure:"udp"` // 新增UDP配置
}

// MySQL配置
type MySQLConfig struct {
	DSN          string `mapstructure:"dsn"`
	MaxOpenConns int    `mapstructure:"max_open_conns" default:"100"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" default:"20"`
}

// Redis配置
type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"DB"`
	PoolSize     int    `mapstructure:"poolSize"`
	MinIdleConns int    `mapstructure:"minIdleConn"`
}

// OSS配置
type OSSConfig struct {
	Endpoint        string `mapstructure:"Endpoint"`
	AccessKeyId     string `mapstructure:"AccessKeyId"`
	AccessKeySecret string `mapstructure:"AccessKeySecret"`
	Bucket          string `mapstructure:"Bucket"`
}

// 超时配置
type TimeoutConfig struct {
	DelayHeartbeat   int           `mapstructure:"DelayHeartbeat"`
	HeartbeatHz      int           `mapstructure:"HeartbeatHz"`
	HeartbeatMaxTime time.Duration `mapstructure:"HeartbeatMaxTime" default:"3000s"` // 转换为Duration
	RedisOnlineTime  int           `mapstructure:"RedisOnlineTime"`
}

// 端口配置
type PortConfig struct {
	Server string `mapstructure:"server"`
	UDP    int    `mapstructure:"udp"`
}

// 应用配置
type AppConfig struct {
	Env      string `mapstructure:"env" default:"dev"` // 环境：dev/test/prod
	Name     string `mapstructure:"name" default:"ginchat"`
	LogLevel string `mapstructure:"log_level" default:"info"`
}

// JWT配置（预留）
type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expire time.Duration `mapstructure:"expire" default:"24h"`
}

// UDP配置（新增）
type UDPConfig struct {
	TargetIP   string `mapstructure:"target_ip" default:"58.198.176.23"` // 目标服务器IP
	TargetPort int    `mapstructure:"target_port" default:"3000"`        // 目标服务器端口
}

var GlobalConfig Config
