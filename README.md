# Ocean Marketing - Gin框架项目骨架

基于Gin框架的完整项目骨架，集成了常用的功能组件，致力于快速的业务研发。

## 🚀 特性

### 核心功能

- ✅ **Gin Web框架** - 高性能HTTP Web框架
- ✅ **配置管理** - 基于Viper的配置文件解析
- ✅ **日志系统** - 基于Zap的结构化日志
- ✅ **数据库支持** - Gorm ORM，支持MySQL/PostgreSQL
- ✅ **Redis缓存** - go-redis客户端封装
- ✅ **JWT认证** - 完整的认证授权系统

### 中间件功能

- ✅ **验证中间件** - 基于govalidator的自动参数验证
- ✅ **接口限流** - 基于内存的限流中间件
- ✅ **跨域支持** - CORS中间件
- ✅ **异常恢复** - Panic恢复 + 飞书通知
- ✅ **链路追踪** - 基于Jaeger的分布式追踪
- ✅ **性能监控** - Prometheus指标收集
- ✅ **性能分析** - pprof性能剖析

### 业务功能

- ✅ **错误码管理** - 统一错误码定义
- ✅ **响应规范** - RESTful API响应格式
- ✅ **邮件发送** - SMTP邮件发送封装
- ✅ **消息队列** - RabbitMQ消息队列支持
- ✅ **类型转换** - 安全的类型转换工具
- ✅ **Swagger文档** - 自动生成API文档

## 📁 项目结构

```
ocean-marketing/
├── cmd/                    # 主要应用程序目录
│   └── server/            # 服务器主程序
│       └── main.go        # 应用入口
├── internal/              # 私有应用程序和库代码
│   ├── config/           # 配置管理
│   ├── handler/          # 控制器层（按模块组织）
│   │   ├── example.go    # 示例控制器
│   │   ├── health.go     # 健康检查控制器
│   │   └── v1/           # V1版本控制器
│   ├── service/          # 业务逻辑层
│   │   └── example.go    # 示例服务
│   ├── model/            # 数据模型
│   │   └── example.go    # 示例模型
│   ├── middleware/       # 中间件
│   │   └── validation.go # 验证中间件
│   ├── router/           # 路由定义
│   └── pkg/             # 内部包
│       ├── database/    # 数据库连接
│       ├── logger/      # 日志系统
│       ├── redis/       # Redis连接
│       └── tracer/      # 链路追踪
├── pkg/                   # 可以被外部应用程序使用的库代码
│   ├── cast/             # 类型转换
│   ├── email/            # 邮件发送
│   ├── errno/            # 错误码定义
│   ├── jwt/              # JWT认证
│   ├── mq/               # 消息队列
│   └── response/         # 响应处理
├── configs/               # 配置文件
│   └── app.yaml          # 应用配置
├── logs/                  # 日志文件目录
├── docs/                  # 文档目录
├── go.mod                 # Go模块文件
└── README.md             # 项目说明
```

## 🛠️ 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL/PostgreSQL
- Redis
- RabbitMQ (可选)
- Jaeger (可选)

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置文件

复制并修改配置文件：

```bash
cp configs/app.yaml.example configs/app.yaml
# 编辑配置文件，设置数据库、Redis等连接信息
```

### 4. 启动服务

```bash
go run cmd/server/main.go
```

### 5. 访问服务

- 应用地址: http://localhost:8080
- 健康检查: http://localhost:8080/health
- API文档: http://localhost:8080/swagger/index.html
- 性能指标: http://localhost:8080/metrics
- 性能分析: http://localhost:8080/debug/pprof/

## 📖 API接口

### Example模块（示例接口）

- `GET /api/v1/examples` - 获取示例列表
- `GET /api/v1/examples/:id` - 获取示例详情
- `POST /api/v1/examples` - 创建示例（需要认证）
- `PUT /api/v1/examples/:id` - 更新示例（需要认证）
- `DELETE /api/v1/examples/:id` - 删除示例（需要认证）

### 系统接口

- `GET /health` - 健康检查
- `GET /ready` - 就绪检查
- `GET /live` - 存活检查

## 🔧 核心功能使用

### 验证中间件

项目支持在路由中直接使用结构体进行参数验证：

```go
// 定义请求结构
type CreateExampleRequest struct {
    Title   string `json:"title" valid:"required,length(1|100)" form:"title"`
    Content string `json:"content" valid:"required,length(1|1000)" form:"content"`
}

// 自定义验证方法
func (r *CreateExampleRequest) Validate() error {
    // 自定义验证逻辑
    return nil
}

// 在路由中使用
exampleGroup.POST("", 
    middleware.Validation(&request.CreateExampleRequest{}), 
    exampleHandler.Create)

// 在控制器中获取验证后的数据
func (h *ExampleHandler) Create(c *gin.Context) {
    var req CreateExampleRequest
    if err := middleware.GetValidatedData(c, &req); err != nil {
        response.Error(c, errno.ErrBind)
        return
    }
    // 使用验证后的数据...
}
```

### 配置管理
```go
import "ocean-marketing/internal/config"

cfg := config.Get()
fmt.Println(cfg.App.Name)
```

### 日志使用
```go
import "ocean-marketing/internal/pkg/logger"
import "go.uber.org/zap"

logger.Info("示例操作", zap.String("action", "create"))
logger.Error("操作失败", zap.Error(err))
```

### 数据库操作
```go
import "ocean-marketing/internal/pkg/database"

db := database.GetDB()
var example Example
db.First(&example, 1)
```

### Redis操作
```go
import "ocean-marketing/internal/pkg/redis"

redis.Set(ctx, "key", "value", time.Hour)
value, err := redis.Get(ctx, "key")
```

### 错误处理
```go
import "ocean-marketing/pkg/errno"
import "ocean-marketing/pkg/response"

// 返回错误
response.Error(c, errno.ErrResourceNotFound)

// 返回成功
response.Success(c, data)
```

## 🏗️ 开发新模块

### 1. 创建模型
```go
// internal/model/product.go
type Product struct {
    ID          uint      `json:"id" gorm:"primarykey"`
    Name        string    `json:"name" gorm:"size:100;not null"`
    Description string    `json:"description" gorm:"type:text"`
    Price       float64   `json:"price" gorm:"type:decimal(10,2)"`
    CreatedBy   string    `json:"created_by" gorm:"size:100"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 2. 创建服务
```go
// internal/service/product.go
type ProductService struct{}

func (s *ProductService) Create(req *CreateProductRequest) (*Product, error) {
    // 业务逻辑
}
```

### 3. 创建控制器
```go
// internal/handler/product.go
type ProductHandler struct {
    productService *ProductService
}

func (h *ProductHandler) Create(c *gin.Context) {
    var req CreateProductRequest
    if err := middleware.GetValidatedData(c, &req); err != nil {
        response.Error(c, errno.ErrBind)
        return
    }
    // 业务处理...
}
```

### 4. 注册路由
```go
// 在router中添加
func setupProductRoutes(g *gin.RouterGroup) {
    productHandler := handler.NewProductHandler()
    
    productGroup := g.Group("/products")
    {
        productGroup.POST("", 
            middleware.Validation(&request.CreateProductRequest{}), 
            productHandler.Create)
    }
}
```

## 🔧 开发工具

### 生成Swagger文档
```bash
# 安装swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/server/main.go
```

### Make命令
```bash
make build          # 编译项目
make run             # 运行项目
make test            # 运行测试
make lint            # 代码检查
make docker-build    # 构建Docker镜像
```

## 📊 监控告警

### Prometheus指标
- HTTP请求总数、延迟、状态码分布
- 活跃连接数
- 请求/响应大小分布

### 飞书告警
当发生panic异常时，自动发送飞书通知，包含：
- 错误信息和堆栈
- 请求路径和IP
- 时间戳

### 链路追踪
集成Jaeger进行分布式链路追踪，可以追踪请求在微服务间的调用链路。

## 📚 文档

- [验证中间件使用指南](docs/validation-middleware-guide.md) - 验证中间件详细使用说明
- [阿里云部署指南](docs/aliyun-deployment.md) - 生产环境部署指南

## 🚦 最佳实践

### 目录规范
- `internal/handler/` - 控制器按模块组织，一个模块一个文件
- `internal/service/` - 业务逻辑层
- `internal/model/` - 数据模型定义
- `pkg/` - 可复用的公共包

### 错误处理
- 使用统一的错误码 (`pkg/errno/`)
- 记录详细的错误日志
- 返回用户友好的错误信息

### 验证规范
- 使用验证中间件进行参数验证
- 在请求结构体中定义验证规则
- 支持自定义验证方法

### 安全规范
- JWT密钥定期轮换
- 敏感配置使用环境变量
- 输入验证和XSS防护

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🔗 相关链接

- [Gin文档](https://gin-gonic.com/)
- [Gorm文档](https://gorm.io/)
- [Viper文档](https://github.com/spf13/viper)
- [Zap文档](https://github.com/uber-go/zap)
- [govalidator文档](https://github.com/asaskevich/govalidator)
 