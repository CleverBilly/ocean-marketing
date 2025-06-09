package model

import (
	"time"

	"gorm.io/gorm"
)

// Example 示例模型
type Example struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Title       string `json:"title" gorm:"size:255;not null" binding:"required,max=255"`
	Description string `json:"description" gorm:"type:text"`
	Status      int    `json:"status" gorm:"default:1;comment:状态 1启用 0禁用"`
	Sort        int    `json:"sort" gorm:"default:0;comment:排序"`
	CreatedBy   string `json:"created_by" gorm:"size:100;comment:创建者"`
}

// TableName 指定表名
func (Example) TableName() string {
	return "examples"
}

// ExampleCreateRequest 创建示例请求
type ExampleCreateRequest struct {
	Title       string `json:"title" binding:"required,max=255"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Sort        int    `json:"sort"`
}

// ExampleUpdateRequest 更新示例请求
type ExampleUpdateRequest struct {
	Title       string `json:"title" binding:"omitempty,max=255"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
	Sort        *int   `json:"sort"`
}

// ExampleResponse 示例响应
type ExampleResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	Sort        int       `json:"sort"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
