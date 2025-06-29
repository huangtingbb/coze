package utils

import (
	"coze-agent-platform/config"
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	cfg := config.Cfg.Log

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// 设置日志格式
	if cfg.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// 设置输出到控制台
	logrus.SetOutput(os.Stdout)
}
