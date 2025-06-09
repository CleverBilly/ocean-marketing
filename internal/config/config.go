package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Email    EmailConfig    `mapstructure:"email"`
	Tracer   TracerConfig   `mapstructure:"tracer"`
	Feishu   FeishuConfig   `mapstructure:"feishu"`
	MQ       MQConfig       `mapstructure:"mq"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	// 阿里云RDS相关配置
	SSLMode      string `mapstructure:"ssl_mode"`
	Timeout      int    `mapstructure:"timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	Loc          string `mapstructure:"loc"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputPath string `mapstructure:"output_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
	Issuer     string `mapstructure:"issuer"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// TracerConfig 链路追踪配置
type TracerConfig struct {
	ServiceName string  `mapstructure:"service_name"`
	AgentHost   string  `mapstructure:"agent_host"`
	AgentPort   int     `mapstructure:"agent_port"`
	SampleRate  float64 `mapstructure:"sample_rate"`
}

// FeishuConfig 飞书配置
type FeishuConfig struct {
	WebhookURL string `mapstructure:"webhook_url"`
}

// MQConfig 消息队列配置
type MQConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Vhost    string `mapstructure:"vhost"`
}

var cfg *Config

// Init 初始化配置
func Init() *Config {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 支持环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到，使用默认配置
		} else {
			panic(err)
		}
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}

	return cfg
}

// Get 获取配置
func Get() *Config {
	return cfg
}

// setDefaults 设置默认配置
func setDefaults() {
	// App默认配置
	viper.SetDefault("app.name", "ocean-marketing")
	viper.SetDefault("app.port", ":8080")
	viper.SetDefault("app.mode", "debug")

	// Database默认配置
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.database", "ocean_marketing")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)

	// Redis默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Log默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output_path", "./logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.max_backups", 10)
	viper.SetDefault("log.compress", true)

	// JWT默认配置
	viper.SetDefault("jwt.secret", "ocean-marketing-secret")
	viper.SetDefault("jwt.expire_time", 3600)
	viper.SetDefault("jwt.issuer", "ocean-marketing")

	// Email默认配置
	viper.SetDefault("email.host", "smtp.gmail.com")
	viper.SetDefault("email.port", 587)

	// Tracer默认配置
	viper.SetDefault("tracer.service_name", "ocean-marketing")
	viper.SetDefault("tracer.agent_host", "localhost")
	viper.SetDefault("tracer.agent_port", 6831)
	viper.SetDefault("tracer.sample_rate", 1.0)

	// MQ默认配置
	viper.SetDefault("mq.driver", "rabbitmq")
	viper.SetDefault("mq.host", "localhost")
	viper.SetDefault("mq.port", 5672)
	viper.SetDefault("mq.username", "guest")
	viper.SetDefault("mq.password", "guest")
	viper.SetDefault("mq.vhost", "/")
}
