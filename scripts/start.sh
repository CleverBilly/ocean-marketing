#!/bin/bash

# Ocean Marketing 启动脚本

set -e

echo "🚀 启动 Ocean Marketing 服务..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.21+"
    exit 1
fi

# 检查配置文件
if [ ! -f "configs/app.yaml" ]; then
    echo "📝 配置文件不存在，复制示例配置..."
    cp configs/app.example.yaml configs/app.yaml 2>/dev/null || echo "⚠️  请手动创建配置文件 configs/app.yaml"
fi

# 创建必要目录
mkdir -p logs
mkdir -p bin

# 下载依赖
echo "📦 下载依赖..."
go mod download

# 生成 Swagger 文档
if command -v swag &> /dev/null; then
    echo "📚 生成 Swagger 文档..."
    swag init -g cmd/server/main.go -o docs
else
    echo "⚠️  swag 未安装，跳过 Swagger 文档生成"
    echo "💡 可通过以下命令安装：go install github.com/swaggo/swag/cmd/swag@latest"
fi

# 编译项目
echo "🔨 编译项目..."
go build -o bin/ocean-marketing cmd/server/main.go

# 启动服务
echo "✅ 启动服务..."
./bin/ocean-marketing

echo ""
echo "🎉 服务启动成功！"
echo "📍 访问地址："
echo "   - 应用: http://localhost:8080"
echo "   - 健康检查: http://localhost:8080/health"
echo "   - API文档: http://localhost:8080/swagger/index.html"
echo "   - 监控指标: http://localhost:8080/metrics"
echo ""
echo "📋 可用API接口："
echo "   - GET  /api/v1/examples          获取示例列表"
echo "   - GET  /api/v1/examples/:id      获取示例详情"
echo "   - POST /api/v1/examples          创建示例（需要认证）"
echo "   - PUT  /api/v1/examples/:id      更新示例（需要认证）"
echo "   - DELETE /api/v1/examples/:id    删除示例（需要认证）" 