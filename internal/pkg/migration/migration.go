package migration

import (
	"ocean-marketing/internal/model"
	"ocean-marketing/internal/pkg/database"
	"ocean-marketing/internal/pkg/logger"

	"go.uber.org/zap"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	db := database.GetDB()

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&model.Example{},
	)

	if err != nil {
		logger.Error("数据库迁移失败", zap.Error(err))
		return err
	}

	logger.Info("数据库迁移成功")
	return nil
}

// CreateTables 手动创建表（如果需要自定义SQL）
func CreateTables() error {
	db := database.GetDB()

	// 这里可以执行自定义的建表SQL
	// 例如：创建索引、触发器等

	// 为示例表添加额外索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_examples_status ON examples(status)").Error; err != nil {
		logger.Error("创建示例状态索引失败", zap.Error(err))
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_examples_created_by ON examples(created_by)").Error; err != nil {
		logger.Error("创建示例创建者索引失败", zap.Error(err))
		return err
	}

	logger.Info("数据库表创建成功")
	return nil
}

// SeedData 种子数据
func SeedData() error {
	db := database.GetDB()

	// 检查是否已经有示例数据
	var count int64
	db.Model(&model.Example{}).Count(&count)

	if count == 0 {
		// 创建默认示例数据
		examples := []*model.Example{
			{
				Title:       "示例标题1",
				Description: "这是第一个示例的描述",
				Status:      1,
				Sort:        1,
				CreatedBy:   "系统",
			},
			{
				Title:       "示例标题2",
				Description: "这是第二个示例的描述",
				Status:      1,
				Sort:        2,
				CreatedBy:   "系统",
			},
		}

		for _, example := range examples {
			if err := db.Create(example).Error; err != nil {
				logger.Error("创建示例数据失败", zap.Error(err))
				return err
			}
		}

		logger.Info("默认示例数据创建成功")
	}

	return nil
}
