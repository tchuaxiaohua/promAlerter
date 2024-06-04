package config

import (
	"fmt"
)

type App struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// GetHost 拼接IP+端口
func (a *App) GetHost() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

type Logger struct {
	Level      string `yaml:"level"`
	FileName   string `yaml:"fileName"`
	MaxSize    int    `yaml:"maxSize"`
	MaxAge     int    `yaml:"maxAge"`
	MaxBackups int    `yaml:"maxBackups"`
	Compress   bool   `yaml:"comPress"`
}

type Config struct {
	App      *App           `yaml:"app"`
	Log      *Logger        `yaml:"log"`
	Channels ChannelConfigs `yaml:"channels"`
	Jvm      *Jvm           `yaml:"jvm"`
}

// NewApp app 默认配置参数
func NewApp() *App {
	return &App{
		Host: "127.0.0.1",
		Port: "8080",
	}
}

// NewConfig Config 默认参数
func NewConfig() *Config {
	return &Config{
		App: NewApp(),
		Log: NewLogger(),
	}
}

// NewLogger 日志默认配置
func NewLogger() *Logger {
	return &Logger{
		Level:      "DEBUG",
		FileName:   "./data/logs/app.log",
		MaxSize:    10,
		MaxAge:     7,
		MaxBackups: 10,
		Compress:   false,
	}
}

// Jvm 导出配置 结构体对象
type Jvm struct {
	DumpMin   int  `yaml:"dump_min"`
	DumpMax   int  `yaml:"dump_max"`
	DumpTsMin int  `yaml:"dump_ts_min"`
	DumpTsMax int  `yaml:"dump_ts_max"`
	IsDump    bool `yaml:"is_dump"`
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	Enabled bool `yaml:"enabled"`
	Configs []struct {
		Name   string `yaml:"name"`
		Token  string `yaml:"token"`
		Secret string `yaml:"secret"`
	} `yaml:"configs"`
}

// FeiShuConfig 飞书配置
type FeiShuConfig struct {
	Enabled bool `yaml:"enabled"`
	Configs []struct {
		Name   string `yaml:"name"`
		Token  string `yaml:"token"`
		Secret string `yaml:"secret"`
	} `yaml:"configs"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	Enabled bool `yaml:"enabled"`
	Configs []struct {
		SMTPServer  string `yaml:"smtpServer"`
		SMTPPort    int    `yaml:"smtpPort"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		FromAddress string `yaml:"fromAddress"`
		ToAddress   string `yaml:"toAddress"`
	} `yaml:"configs"`
}

// ChannelConfigs  告警发送渠道结构体
type ChannelConfigs struct {
	Email    EmailConfig    `yaml:"email"`
	DingTalk DingTalkConfig `yaml:"dingtalk"`
	FeiShu   FeiShuConfig   `yaml:"feishu"`
}
