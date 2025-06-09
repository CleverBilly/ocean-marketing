# 验证中间件使用指南

## 概述

项目集成了基于govalidator的验证中间件，支持在路由中直接使用结构体进行参数验证。

## 核心特性

1. **自动验证** - 请求参数在到达控制器前已经验证完成
2. **类型安全** - 通过反射确保数据类型正确
3. **灵活组织** - 可以在路由文件中直接定义控制器
4. **上下文传递** - 验证后的数据通过gin.Context传递

## 验证中间件使用

### 支持的验证类型

1. **JSON验证** - `middleware.Validation()`
2. **查询参数验证** - `middleware.ValidationQuery()`
3. **表单验证** - `middleware.ValidationForm()`

### 基本用法

```go
// 在路由中使用验证中间件
exampleGroup.POST("", 
    middleware.Validation(&request.CreateExampleRequest{}), 
    exampleHandler.Create)
```

### 在控制器中获取验证数据

```go
func (h *ExampleHandler) Create(c *gin.Context) {
    // 从上下文获取验证后的数据
    var req CreateExampleRequest
    if err := middleware.GetValidatedData(c, &req); err != nil {
        response.Error(c, errno.ErrBind)
        return
    }
    
    // 使用验证后的数据进行业务处理...
}
```

## 请求结构定义

### 基本结构

```go
// 创建示例请求
type CreateExampleRequest struct {
    Title   string `json:"title" valid:"required,length(1|100)" form:"title"`
    Content string `json:"content" valid:"required,length(1|1000)" form:"content"`
    Status  int    `json:"status" valid:"optional,range(0|1)" form:"status"`
}

// 自定义验证方法
func (r *CreateExampleRequest) Validate() error {
    // 自定义验证逻辑
    if r.Status < 0 || r.Status > 1 {
        return errors.New("状态值只能是0或1")
    }
    return nil
}
```

### 验证标签说明

- `required` - 必填字段
- `optional` - 可选字段
- `email` - 邮箱格式
- `length(min|max)` - 长度限制
- `range(min|max)` - 数值范围
- `url` - URL格式

## 路由组织方式

### 推荐的组织模式

```go
func setupExampleRoutes(g *gin.RouterGroup) {
    // 创建控制器
    exampleHandler := handler.NewExampleHandler()

    // 示例路由组
    exampleGroup := g.Group("/examples")
    {
        // 获取列表（查询参数验证）
        exampleGroup.GET("", 
            middleware.ValidationQuery(&request.ExampleListRequest{}), 
            exampleHandler.GetList)
        
        // 创建示例（JSON验证）
        exampleGroup.POST("", 
            middleware.AuthMiddleware(),
            middleware.Validation(&request.CreateExampleRequest{}), 
            exampleHandler.Create)
        
        // 更新示例（JSON验证）
        exampleGroup.PUT("/:id", 
            middleware.AuthMiddleware(),
            middleware.Validation(&request.UpdateExampleRequest{}), 
            exampleHandler.Update)
    }
}
```

### 灵活的写法风格

```go
func exampleApi(g *gin.RouterGroup) {
    exampleHandler := handler.NewExampleHandler()

    exampleGroup := g.Group("/examples")
    {
        // 获取示例列表
        exampleGroup.GET("", 
            middleware.ValidationQuery(&request.ExampleListRequest{}), 
            exampleHandler.GetList)
        
        // 创建示例
        exampleGroup.POST("", 
            middleware.Validation(&request.CreateExampleRequest{}), 
            exampleHandler.Create)
    }
}
```

## 完整示例

### 1. 定义请求结构

```go
// internal/request/example_request.go (如果使用request包)
type CreateExampleRequest struct {
    Title       string `json:"title" valid:"required,length(1|100)" form:"title"`
    Description string `json:"description" valid:"optional,length(0|500)" form:"description"`
    Status      int    `json:"status" valid:"optional,range(0|1)" form:"status"`
}

func (r *CreateExampleRequest) Validate() error {
    if r.Status < 0 || r.Status > 1 {
        return errors.New("状态值只能是0或1")
    }
    return nil
}
```

### 2. 控制器处理

```go
// internal/handler/example.go
func (h *ExampleHandler) Create(c *gin.Context) {
    // 获取验证后的数据
    var req CreateExampleRequest
    if err := middleware.GetValidatedData(c, &req); err != nil {
        zap.L().Error("获取验证数据失败", zap.Error(err))
        response.Error(c, errno.ErrBind)
        return
    }

    // 业务逻辑处理
    example, err := h.exampleService.Create(&req, "admin")
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, example)
}
```

### 3. 路由注册

```go
// internal/router/router.go
func Register(r *gin.Engine, cfg *config.Config) {
    api := r.Group("/api")
    {
        v1Group := api.Group("/v1")
        setupExampleRoutes(v1Group)
    }
}
```

## API使用示例

### 创建示例

```bash
POST /api/v1/examples
Content-Type: application/json

{
    "title": "示例标题",
    "description": "示例描述",
    "status": 1
}
```

### 获取示例列表

```bash
GET /api/v1/examples?page=1&size=10&status=1
```

### 更新示例

```bash
PUT /api/v1/examples/1
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "title": "更新后的标题",
    "description": "更新后的描述"
}
```

## 验证错误处理

当验证失败时，中间件会自动返回错误响应：

```json
{
    "code": 10002,
    "message": "请求参数错误",
    "data": null
}
```

如果有自定义验证错误：

```json
{
    "code": 10003,
    "message": "状态值只能是0或1",
    "data": null
}
```

## 最佳实践

1. **验证规则优先级**
   - 优先使用govalidator标签进行基础验证
   - 复杂业务逻辑使用`Validate()`方法

2. **错误处理**
   - 使用统一的错误码定义
   - 返回用户友好的错误信息

3. **性能考虑**
   - 验证在中间件层完成，避免重复验证
   - 使用反射缓存提高性能

4. **代码组织**
   - 请求结构体定义在合适的包中
   - 验证逻辑与业务逻辑分离

