package database

import (
	"fmt"
	"time"

	"ocean-marketing/internal/config"
	"ocean-marketing/internal/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init(cfg config.DatabaseConfig) {
	var err error
	var dsn string

	switch cfg.Driver {
	case "mysql":
		// 构建基础DSN
		loc := cfg.Loc
		if loc == "" {
			loc = "Local"
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset, loc)

		// 添加SSL配置（阿里云RDS推荐）
		if cfg.SSLMode != "" {
			dsn += "&tls=" + cfg.SSLMode
		}

		// 添加超时配置
		if cfg.Timeout > 0 {
			dsn += fmt.Sprintf("&timeout=%ds", cfg.Timeout)
		}
		if cfg.ReadTimeout > 0 {
			dsn += fmt.Sprintf("&readTimeout=%ds", cfg.ReadTimeout)
		}
		if cfg.WriteTimeout > 0 {
			dsn += fmt.Sprintf("&writeTimeout=%ds", cfg.WriteTimeout)
		}

		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		})
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		})
	default:
		logger.Fatal("不支持的数据库驱动", zap.String("driver", cfg.Driver))
	}

	if err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}

	// 获取底层的sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Fatal("获取数据库实例失败", zap.Error(err))
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		logger.Fatal("数据库连接测试失败", zap.Error(err))
	}

	logger.Info("数据库连接成功", zap.String("driver", cfg.Driver))
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
