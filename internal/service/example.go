package service

import (
	"errors"

	"ocean-marketing/internal/model"
	"ocean-marketing/internal/pkg/database"
	"ocean-marketing/pkg/errno"

	"gorm.io/gorm"
)

// ExampleService 示例服务
type ExampleService struct{}

// NewExampleService 创建示例服务实例
func NewExampleService() *ExampleService {
	return &ExampleService{}
}

// GetList 获取示例列表
func (s *ExampleService) GetList(page, size int) ([]model.ExampleResponse, int64, error) {
	var examples []model.Example
	var total int64

	db := database.GetDB().Model(&model.Example{})

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, errno.ErrDatabase
	}

	// 分页查询
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Find(&examples).Error; err != nil {
		return nil, 0, errno.ErrDatabase
	}

	// 转换为响应格式
	var responses []model.ExampleResponse
	for _, example := range examples {
		response := model.ExampleResponse{
			ID:          example.ID,
			Title:       example.Title,
			Description: example.Description,
			Status:      example.Status,
			Sort:        example.Sort,
			CreatedBy:   example.CreatedBy,
			CreatedAt:   example.CreatedAt,
			UpdatedAt:   example.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetByID 根据ID获取示例
func (s *ExampleService) GetByID(id uint) (*model.ExampleResponse, error) {
	var example model.Example
	if err := database.GetDB().First(&example, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrResourceNotFound
		}
		return nil, errno.ErrDatabase
	}

	response := &model.ExampleResponse{
		ID:          example.ID,
		Title:       example.Title,
		Description: example.Description,
		Status:      example.Status,
		Sort:        example.Sort,
		CreatedBy:   example.CreatedBy,
		CreatedAt:   example.CreatedAt,
		UpdatedAt:   example.UpdatedAt,
	}

	return response, nil
}

// Create 创建示例
func (s *ExampleService) Create(req *model.ExampleCreateRequest, createdBy string) (*model.ExampleResponse, error) {
	example := &model.Example{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Sort:        req.Sort,
		CreatedBy:   createdBy,
	}

	if err := database.GetDB().Create(example).Error; err != nil {
		return nil, errno.ErrDatabase
	}

	return s.GetByID(example.ID)
}

// Update 更新示例
func (s *ExampleService) Update(id uint, req *model.ExampleUpdateRequest, currentUser string) (*model.ExampleResponse, error) {
	var example model.Example
	if err := database.GetDB().First(&example, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrResourceNotFound
		}
		return nil, errno.ErrDatabase
	}

	// 检查权限（只有创建者可以修改）
	if example.CreatedBy != currentUser && currentUser != "admin" {
		return nil, errno.ErrPermissionDenied
	}

	// 更新字段
	if req.Title != "" {
		example.Title = req.Title
	}
	if req.Description != "" {
		example.Description = req.Description
	}
	if req.Status != nil {
		example.Status = *req.Status
	}
	if req.Sort != nil {
		example.Sort = *req.Sort
	}

	if err := database.GetDB().Save(&example).Error; err != nil {
		return nil, errno.ErrDatabase
	}

	return s.GetByID(example.ID)
}

// Delete 删除示例
func (s *ExampleService) Delete(id uint, currentUser string) error {
	var example model.Example
	if err := database.GetDB().First(&example, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrResourceNotFound
		}
		return errno.ErrDatabase
	}

	// 检查权限（只有创建者可以删除）
	if example.CreatedBy != currentUser && currentUser != "admin" {
		return errno.ErrPermissionDenied
	}

	if err := database.GetDB().Delete(&example).Error; err != nil {
		return errno.ErrDatabase
	}

	return nil
}
