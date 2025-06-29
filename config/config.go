package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Coze     CozeConfig     `mapstructure:"coze"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type CozeConfig struct {
	APIURL             string `mapstructure:"api_url"`
	APIKey             string `mapstructure:"api_key"`
	ClientID           string `mapstructure:"client_id"`
	PublicKeyID        string `mapstructure:"public_key_id"`
	PrivateKey         string `mapstructure:"private_key"`
	PrivateKeyFilePath string `mapstructure:"private_key_file_path"`
}

var Cfg *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")    // 支持从cmd目录启动时查找配置
	viper.AddConfigPath("../../config") // 支持更深层次的目录

	// 设置环境变量
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("配置文件读取失败，使用默认配置: %v", err)
	} else {
		log.Printf("成功读取配置文件: %s", viper.ConfigFileUsed())
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		log.Fatal("配置解析失败:", err)
	}

	log.Printf("配置初始化完成, Coze配置: %+v", Cfg.Coze)
}

// Get 获取全局配置
func Get() *Config {
	return Cfg
}

// GetCozeConfig 获取 Coze 配置
func GetCozeConfig() *CozeConfig {
	return &Cfg.Coze
}

func setDefaults() {
	// App配置默认值
	viper.SetDefault("app.name", "coze-agent-platform")
	viper.SetDefault("app.port", "8080")
	viper.SetDefault("app.mode", "debug")

	// 数据库配置默认值
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "coze_agent")

	// Redis配置默认值
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// JWT配置默认值
	viper.SetDefault("jwt.secret", "coze-agent-secret-key")
	viper.SetDefault("jwt.expire", 7200) // 2小时

	// 日志配置默认值
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// Coze配置默认值
	viper.SetDefault("coze.api_key", "")
	viper.SetDefault("coze.api_url", "https://www.coze.cn/api")
	viper.SetDefault("coze.client_id", "")
	viper.SetDefault("coze.public_key_id", "")
	viper.SetDefault("coze.private_key", "")
	viper.SetDefault("coze.private_key_file_path", "")
}
