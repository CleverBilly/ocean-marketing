# Ocean Marketing 项目概述

## 项目状态

当前项目已完成架构，保留了核心框架和基础功能，同时提供了完整的Example示例模块作为开发参考。

## 当前功能模块

### ✅ 核心基础设施
- **配置管理** - 基于Viper的配置系统
- **日志系统** - 基于Zap的结构化日志
- **数据库** - Gorm ORM支持
- **Redis缓存** - go-redis客户端
- **JWT认证** - 完整的JWT认证系统

### ✅ 中间件系统
- **验证中间件** - 基于govalidator的自动参数验证
- **认证中间件** - JWT令牌验证
- **限流中间件** - 接口访问限流
- **CORS中间件** - 跨域请求支持
- **异常恢复** - Panic恢复和飞书通知
- **链路追踪** - Jaeger分布式追踪
- **性能监控** - Prometheus指标收集

### ✅ 业务模块
- **Example模块** - 完整的CRUD示例
  - 数据模型定义 (`internal/model/example.go`)
  - 业务逻辑层 (`internal/service/example.go`)
  - 控制器层 (`internal/handler/example.go`)
  - V1版本控制器 (`internal/handler/v1/example.go`)

### ✅ 系统接口
- **健康检查** - `/health`, `/ready`, `/live`
- **监控接口** - `/metrics` (Prometheus指标)
- **性能分析** - `/debug/pprof/` (性能剖析)
- **API文档** - `/swagger/index.html` (Swagger文档)

## 可用API接口

### Example模块 (示例接口)
- `GET /api/v1/examples` - 获取示例列表
- `GET /api/v1/examples/:id` - 获取示例详情  
- `POST /api/v1/examples` - 创建示例 (需要认证)
- `PUT /api/v1/examples/:id` - 更新示例 (需要认证)
- `DELETE /api/v1/examples/:id` - 删除示例 (需要认证)

### 系统接口
- `GET /health` - 健康检查
- `GET /ready` - 就绪检查
- `GET /live` - 存活检查


## 开发指南

### 快速开始
1. 配置数据库和Redis连接信息
2. 运行 `go run cmd/server/main.go` 启动服务
3. 访问 `http://localhost:8080/health` 验证服务状态
4. 访问 `http://localhost:8080/swagger/index.html` 查看API文档

### 开发新模块
1. 参考Example模块的实现
2. 按照相同的架构模式组织代码
3. 使用验证中间件进行参数验证
4. 遵循统一的错误处理和响应格式

### 架构特点
- **分层架构** - Handler/Service/Model分离
- **依赖注入** - 清晰的依赖关系管理
- **中间件驱动** - 功能通过中间件组合
- **验证优先** - 自动参数验证
- **类型安全** - 强类型定义和检查

## 文档资源

- [项目README](../README.md) - 项目完整介绍和使用指南
- [验证中间件指南](validation-middleware-guide.md) - 验证中间件详细使用说明
- [阿里云部署指南](aliyun-deployment.md) - 生产环境部署说明

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: Gorm
- **数据库**: MySQL/PostgreSQL
- **缓存**: Redis
- **配置**: Viper
- **日志**: Zap
- **认证**: JWT
- **验证**: govalidator
- **监控**: Prometheus + Jaeger
- **容器**: Docker

## 贡献指南

1. 基于Example模块创建新功能
2. 遵循现有的代码风格和架构模式
3. 添加必要的单元测试
4. 更新相关文档

项目现在处于一个干净的状态，可以作为新项目的起点或继续开发新的业务模块。 