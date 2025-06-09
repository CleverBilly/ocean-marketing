package handler

import (
	"strconv"

	"ocean-marketing/internal/middleware"
	"ocean-marketing/internal/model"
	"ocean-marketing/internal/service"
	"ocean-marketing/pkg/errno"
	"ocean-marketing/pkg/response"

	"github.com/gin-gonic/gin"
)

// ExampleHandler 示例控制器
type ExampleHandler struct {
	exampleService *service.ExampleService
}

// NewExampleHandler 创建示例控制器实例
func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{
		exampleService: service.NewExampleService(),
	}
}

// GetExamples 获取示例列表
// @Summary 获取示例列表
// @Description 分页获取示例列表
// @Tags 示例管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageResponse} "获取成功"
// @Router /api/v1/examples [get]
func (h *ExampleHandler) GetExamples(c *gin.Context) {
	// 获取分页参数
	page := 1
	size := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := c.Query("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			size = s
		}
	}

	list, total, err := h.exampleService.GetList(page, size)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, list, total, page, size)
}

// GetExample 获取单个示例
// @Summary 获取单个示例
// @Description 根据ID获取示例详情
// @Tags 示例管理
// @Accept json
// @Produce json
// @Param id path int true "示例ID"
// @Success 200 {object} response.Response{data=model.ExampleResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "示例不存在"
// @Router /api/v1/examples/{id} [get]
func (h *ExampleHandler) GetExample(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, errno.ErrBind)
		return
	}

	example, err := h.exampleService.GetByID(uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, example)
}

// CreateExample 创建示例
// @Summary 创建示例
// @Description 创建新的示例
// @Tags 示例管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body model.ExampleCreateRequest true "创建信息"
// @Success 200 {object} response.Response{data=model.ExampleResponse} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/v1/examples [post]
func (h *ExampleHandler) CreateExample(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		response.Unauthorized(c, errno.ErrTokenInvalid)
		return
	}

	var req model.ExampleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errno.ErrBind)
		return
	}

	// 创建示例，使用userID作为创建者
	currentUser := strconv.FormatUint(uint64(userID), 10)
	example, err := h.exampleService.Create(&req, currentUser)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, example)
}

// UpdateExample 更新示例
// @Summary 更新示例
// @Description 更新指定ID的示例
// @Tags 示例管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "示例ID"
// @Param request body model.ExampleUpdateRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.ExampleResponse} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "示例不存在"
// @Router /api/v1/examples/{id} [put]
func (h *ExampleHandler) UpdateExample(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		response.Unauthorized(c, errno.ErrTokenInvalid)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, errno.ErrBind)
		return
	}

	var req model.ExampleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errno.ErrBind)
		return
	}

	// 更新示例，使用userID作为当前用户
	currentUser := strconv.FormatUint(uint64(userID), 10)
	example, err := h.exampleService.Update(uint(id), &req, currentUser)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, example)
}

// DeleteExample 删除示例
// @Summary 删除示例
// @Description 删除指定ID的示例
// @Tags 示例管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "示例ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "示例不存在"
// @Router /api/v1/examples/{id} [delete]
func (h *ExampleHandler) DeleteExample(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		response.Unauthorized(c, errno.ErrTokenInvalid)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, errno.ErrBind)
		return
	}

	// 删除示例，使用userID作为当前用户
	currentUser := strconv.FormatUint(uint64(userID), 10)
	if err := h.exampleService.Delete(uint(id), currentUser); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
